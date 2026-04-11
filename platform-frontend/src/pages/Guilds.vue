<script setup lang="ts">
import type { GuildMemberEntry, GuildPublic } from '@/models/guild'
import { Loader2, LogOut, Plus, Shield, Trash2, Users } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import GuildCardSkeleton from '@/components/guilds/GuildCardSkeleton.vue'
import {
  Dialog,
  DialogFooter,
  DialogHeader,
  DialogScrollContent,
  DialogTitle,
} from '@/components/ui/dialog'
import { Typography } from '@/components/ui/typography'
import { useSSE } from '@/composables/useSSE'
import { useUser } from '@/composables/useUser'
import { displayName } from '@/lib/utils'
import { handleError } from '@/services/errorService'
import { guildService } from '@/services/guilds'

const guilds = ref<GuildPublic[]>([])
const isLoading = ref(true)
const isSubmitting = ref(false)
const showCreateDialog = ref(false)
const showMembersDialog = ref(false)
const selectedGuildMembers = ref<GuildMemberEntry[]>([])
const selectedGuildName = ref('')

const user = useUser()
const actionInProgress = ref<number | null>(null)

const newName = ref('')
const newDescription = ref('')
const newIcon = ref('users')
const newColor = ref('#6366f1')

const colorOptions = ['#6366f1', '#ec4899', '#f59e0b', '#10b981', '#3b82f6', '#8b5cf6', '#ef4444', '#06b6d4']

async function fetchGuilds() {
  isLoading.value = true
  try {
    guilds.value = await guildService.getAll() ?? []
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

async function createGuild() {
  if (!newName.value.trim())
    return
  isSubmitting.value = true
  try {
    await guildService.create({
      name: newName.value.trim(),
      description: newDescription.value.trim(),
      icon: newIcon.value,
      color: newColor.value,
    })
    showCreateDialog.value = false
    newName.value = ''
    newDescription.value = ''
    await fetchGuilds()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSubmitting.value = false
  }
}

async function joinGuild(id: number) {
  if (actionInProgress.value)
    return
  actionInProgress.value = id
  try {
    await guildService.join(id)
    await fetchGuilds()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    actionInProgress.value = null
  }
}

async function leaveGuild(id: number) {
  if (actionInProgress.value)
    return
  actionInProgress.value = id
  try {
    await guildService.leave(id)
    await fetchGuilds()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    actionInProgress.value = null
  }
}

async function deleteGuild(id: number) {
  if (actionInProgress.value)
    return
  actionInProgress.value = id
  try {
    await guildService.remove(id)
    await fetchGuilds()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    actionInProgress.value = null
  }
}

async function showMembers(guild: GuildPublic) {
  selectedGuildName.value = guild.name
  try {
    selectedGuildMembers.value = await guildService.getMembers(guild.id)
    showMembersDialog.value = true
  }
  catch (error) {
    handleError(error)
  }
}

function isInAnyGuild() {
  return guilds.value.some(g => g.isMember)
}

useSSE('guilds', () => fetchGuilds())

onMounted(() => {
  fetchGuilds()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="flex items-center justify-between mb-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Гильдии
      </Typography>
      <button
        v-if="!isInAnyGuild()"
        class="flex items-center gap-2 px-4 py-2 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors"
        @click="showCreateDialog = true"
      >
        <Plus class="h-4 w-4" />
        Создать
      </button>
    </div>

    <div
      v-if="isLoading"
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
    >
      <GuildCardSkeleton v-for="i in 6" :key="i" />
    </div>

    <template v-else>
      <div
        v-if="guilds.length === 0"
        class="text-center py-12 text-muted-foreground"
      >
        <Shield class="h-12 w-12 mx-auto mb-3 opacity-50" />
        <p>Гильдий пока нет. Создайте первую!</p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="guild in guilds"
          :key="guild.id"
          class="rounded-2xl border bg-card border-border overflow-hidden"
        >
          <div
            class="h-2"
            :style="{ backgroundColor: guild.color }"
          />
          <div class="p-4">
            <div class="flex items-start justify-between mb-2">
              <h3 class="font-semibold">
                {{ guild.name }}
              </h3>
              <span
                v-if="guild.isMember"
                class="px-2 py-0.5 rounded-full text-xs bg-primary/10 text-primary font-medium"
              >
                Участник
              </span>
            </div>

            <p
              v-if="guild.description"
              class="text-sm text-muted-foreground mb-3"
            >
              {{ guild.description }}
            </p>

            <div class="flex items-center gap-4 text-sm text-muted-foreground mb-3">
              <button
                class="flex items-center gap-1 hover:text-foreground transition-colors"
                @click="showMembers(guild)"
              >
                <Users class="h-3.5 w-3.5" />
                {{ guild.memberCount }} участников
              </button>
              <div class="flex items-center gap-1">
                <span class="font-bold text-foreground">{{ guild.totalPoints }}</span> баллов
              </div>
            </div>

            <div class="flex items-center gap-1.5 text-xs text-muted-foreground mb-3">
              <img
                :src="guild.ownerAvatarUrl || `https://ui-avatars.com/api/?name=${encodeURIComponent(guild.ownerFirstName || '?')}&background=random`"
                :alt="displayName(guild.ownerFirstName, guild.ownerLastName)"
                class="h-5 w-5 rounded-full object-cover"
              >
              <span>{{ displayName(guild.ownerFirstName, guild.ownerLastName) }}</span>
            </div>

            <div class="flex gap-2">
              <button
                v-if="!guild.isMember && !isInAnyGuild()"
                class="flex-1 px-3 py-1.5 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
                :disabled="actionInProgress === guild.id"
                @click="joinGuild(guild.id)"
              >
                Вступить
              </button>
              <button
                v-if="guild.isMember && guild.ownerId !== user?.id"
                class="flex items-center gap-1 px-3 py-1.5 rounded-xl text-sm font-medium text-muted-foreground hover:bg-muted transition-colors disabled:opacity-50"
                :disabled="actionInProgress === guild.id"
                @click="leaveGuild(guild.id)"
              >
                <LogOut class="h-3.5 w-3.5" />
                Выйти
              </button>
              <button
                v-if="guild.ownerId === user?.id"
                class="flex items-center gap-1 px-3 py-1.5 rounded-xl text-sm font-medium text-red-500 hover:bg-red-500/10 transition-colors ml-auto disabled:opacity-50"
                :disabled="actionInProgress === guild.id"
                @click="deleteGuild(guild.id)"
              >
                <Trash2 class="h-3.5 w-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Create dialog -->
    <Dialog v-model:open="showCreateDialog">
      <DialogScrollContent>
        <DialogHeader>
          <DialogTitle>Создать гильдию</DialogTitle>
        </DialogHeader>
        <form
          class="space-y-4"
          @submit.prevent="createGuild"
        >
          <div>
            <label class="block text-sm font-medium mb-1">Название *</label>
            <input
              v-model="newName"
              type="text"
              required
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="Название гильдии"
            >
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">Описание</label>
            <textarea
              v-model="newDescription"
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary min-h-16 resize-none"
              placeholder="О чём ваша гильдия..."
            />
          </div>
          <div>
            <label class="block text-sm font-medium mb-2">Цвет</label>
            <div class="flex gap-2">
              <button
                v-for="color in colorOptions"
                :key="color"
                type="button"
                class="h-8 w-8 rounded-full border-2 transition-transform"
                :class="newColor === color ? 'border-foreground scale-110' : 'border-transparent'"
                :style="{ backgroundColor: color }"
                @click="newColor = color"
              />
            </div>
          </div>
          <DialogFooter>
            <button
              type="submit"
              class="px-4 py-2 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
              :disabled="!newName.trim() || isSubmitting"
            >
              <Loader2
                v-if="isSubmitting"
                class="h-4 w-4 animate-spin inline mr-1"
              />
              Создать
            </button>
          </DialogFooter>
        </form>
      </DialogScrollContent>
    </Dialog>

    <!-- Members dialog -->
    <Dialog v-model:open="showMembersDialog">
      <DialogScrollContent>
        <DialogHeader>
          <DialogTitle>{{ selectedGuildName }}</DialogTitle>
        </DialogHeader>
        <div class="space-y-2">
          <div
            v-for="member in selectedGuildMembers"
            :key="member.memberId"
            class="flex items-center gap-3 p-2 rounded-xl"
          >
            <img
              :src="member.avatarUrl || `https://ui-avatars.com/api/?name=${encodeURIComponent(member.firstName || '?')}&background=random`"
              :alt="displayName(member.firstName, member.lastName)"
              class="h-8 w-8 rounded-full object-cover"
            >
            <router-link
              :to="`/members/${member.memberId}`"
              class="flex-1 text-sm font-medium hover:underline truncate"
            >
              {{ displayName(member.firstName, member.lastName) }}
            </router-link>
            <span class="text-sm font-bold tabular-nums">{{ member.total }} б.</span>
          </div>
        </div>
      </DialogScrollContent>
    </Dialog>
  </div>
</template>
