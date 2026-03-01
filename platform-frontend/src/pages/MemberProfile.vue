<script setup lang="ts">
import type { PublicProfile } from '@/models/profile'
import { Typography } from 'itx-ui-kit'
import { ArrowLeft, Loader2, Trophy } from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { handleError } from '@/services/errorService'
import { profileService } from '@/services/profile'

const route = useRoute()
const router = useRouter()
const profile = ref<PublicProfile | null>(null)
const isLoading = ref(true)

function pluralizeDays(n: number): string {
  if (n % 10 === 1 && n % 100 !== 11)
    return 'день'
  if (n % 10 >= 2 && n % 10 <= 4 && (n % 100 < 10 || n % 100 >= 20))
    return 'дня'
  return 'дней'
}

const daysSinceJoined = computed(() => {
  if (!profile.value?.member.createdAt)
    return null
  const date = new Date(profile.value.member.createdAt)
  if (date.getFullYear() <= 1)
    return null
  const diff = Date.now() - date.getTime()
  return Math.max(1, Math.floor(diff / 86400000))
})

function getAvatarSrc(tg: string, avatarUrl?: string) {
  return avatarUrl || `https://t.me/i/userpic/320/${tg}.jpg`
}

async function loadProfile() {
  isLoading.value = true
  try {
    const id = Number(route.params.id)
    profile.value = await profileService.getMemberById(id)
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

onMounted(loadProfile)
</script>

<template>
  <div class="container mx-auto px-4 py-8">
    <button
      class="flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground mb-6"
      @click="router.back()"
    >
      <ArrowLeft class="h-4 w-4" />
      Назад
    </button>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <div
      v-else-if="profile"
      class="space-y-6"
    >
      <div class="bg-card rounded-3xl border p-6">
        <div class="flex items-start gap-5">
          <div class="w-20 h-20 rounded-full overflow-hidden shrink-0 bg-accent/20">
            <img
              :src="getAvatarSrc(profile.member.tg, profile.member.avatarUrl)"
              class="w-full h-full object-cover"
            >
          </div>
          <div class="min-w-0">
            <Typography
              variant="h2"
              as="h1"
              class="mb-1"
            >
              {{ profile.member.firstName }} {{ profile.member.lastName }}
            </Typography>
            <a
              v-if="profile.member.tg"
              :href="`https://t.me/${profile.member.tg}`"
              target="_blank"
              class="text-sm text-primary underline"
            >
              @{{ profile.member.tg }}
            </a>
          </div>
        </div>

        <p
          v-if="profile.member.bio"
          class="mt-4 text-muted-foreground"
        >
          {{ profile.member.bio }}
        </p>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div class="bg-card rounded-3xl border p-6 flex items-center gap-3">
          <Trophy class="h-5 w-5 text-yellow-500 shrink-0" />
          <div>
            <div class="text-2xl font-bold">
              {{ profile.points }}
            </div>
            <div class="text-sm text-muted-foreground">
              Баллов
            </div>
          </div>
        </div>
        <div
          v-if="daysSinceJoined"
          class="bg-card rounded-3xl border p-6"
        >
          <div class="text-2xl font-bold">
            {{ daysSinceJoined }}
          </div>
          <div class="text-sm text-muted-foreground">
            {{ pluralizeDays(daysSinceJoined) }} с нами
          </div>
        </div>
      </div>

      <div
        v-if="profile.isMentor && profile.mentor"
        class="bg-card rounded-3xl border p-6"
      >
        <Typography
          variant="h3"
          as="h2"
          class="mb-2"
        >
          Ментор
        </Typography>
        <p
          v-if="profile.mentor.occupation"
          class="text-muted-foreground mb-3"
        >
          {{ profile.mentor.occupation }}
        </p>
        <RouterLink
          :to="`/mentors/${profile.mentor.id}`"
          class="text-sm text-primary underline"
        >
          Перейти к профилю ментора
        </RouterLink>
      </div>
    </div>
  </div>
</template>
