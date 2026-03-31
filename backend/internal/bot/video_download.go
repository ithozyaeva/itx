package bot

import (
	"context"
	"fmt"
	"ithozyeva/config"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var videoURLRegex = regexp.MustCompile(
	`https?://(?:` +
		`(?:www\.)?instagram\.com/reels?/[\w-]+` +
		`|(?:www\.)?tiktok\.com/@[\w.]+/video/\d+` +
		`|(?:vm|vt)\.tiktok\.com/[\w-]+` +
		`|(?:www\.)?youtube\.com/shorts/[\w-]+` +
		`|youtu\.be/[\w-]+` +
		`)[\w\-._~:/?#\[\]@!$&'()*+,;=%]*`,
)

const maxVideoURLsPerMessage = 3

func extractVideoURLs(text string) []string {
	matches := videoURLRegex.FindAllString(text, maxVideoURLsPerMessage+1)
	if len(matches) > maxVideoURLsPerMessage {
		matches = matches[:maxVideoURLsPerMessage]
	}
	return matches
}

func (b *TelegramBot) handleVideoURLs(message *tgbotapi.Message, urls []string) {
	if message.From != nil && message.From.ID == b.bot.Self.ID {
		return
	}

	// Работает только в группах ITX (основной чат + tracked chats)
	if message.Chat.ID != config.CFG.TelegramMainChatID && !b.chatActivityService.IsTrackedChat(message.Chat.ID) {
		return
	}

	for _, url := range urls {
		if err := b.downloadAndSendVideo(message.Chat.ID, message.MessageID, url); err != nil {
			log.Printf("[video_download] error processing %s: %v", url, err)
		}
	}
}

func (b *TelegramBot) downloadAndSendVideo(chatID int64, replyToMsgID int, url string) error {
	tmpDir, err := os.MkdirTemp("", "ytdlp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	outputTemplate := filepath.Join(tmpDir, "video.%(ext)s")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "yt-dlp",
		"--no-playlist",
		"--max-filesize", "50m",
		"-f", "best[filesize<50M]/best",
		"-o", outputTemplate,
		url,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			b.sendReplyText(chatID, replyToMsgID, "Скачивание видео заняло слишком много времени")
			return fmt.Errorf("timeout downloading %s", url)
		}
		b.sendReplyText(chatID, replyToMsgID, "Не удалось скачать видео по этой ссылке")
		return fmt.Errorf("yt-dlp failed: %w, output: %s", err, string(output))
	}

	files, err := filepath.Glob(filepath.Join(tmpDir, "video.*"))
	if err != nil || len(files) == 0 {
		b.sendReplyText(chatID, replyToMsgID, "Не удалось скачать видео по этой ссылке")
		return fmt.Errorf("no downloaded file found in %s", tmpDir)
	}

	videoPath := files[0]

	info, err := os.Stat(videoPath)
	if err != nil {
		return fmt.Errorf("failed to stat video file: %w", err)
	}

	const maxSize = 49 * 1024 * 1024 // 49 MB
	if info.Size() > maxSize {
		b.sendReplyText(chatID, replyToMsgID, "Видео слишком большое для Telegram (лимит 50 МБ)")
		return fmt.Errorf("video too large: %d bytes", info.Size())
	}

	video := tgbotapi.NewVideo(chatID, tgbotapi.FilePath(videoPath))
	video.ReplyToMessageID = replyToMsgID

	if _, err := b.bot.Send(video); err != nil {
		return fmt.Errorf("failed to send video: %w", err)
	}

	return nil
}

func (b *TelegramBot) sendReplyText(chatID int64, replyToMsgID int, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyToMessageID = replyToMsgID
	if _, err := b.bot.Send(msg); err != nil {
		log.Printf("[video_download] failed to send reply: %v", err)
	}
}
