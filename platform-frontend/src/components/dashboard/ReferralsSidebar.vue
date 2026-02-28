<script setup lang="ts">
import type { ReferalLink } from '@/models/referals'
import { Loader2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import ReferralCard from '@/components/dashboard/ReferralCard.vue'
import { handleError } from '@/services/errorService'
import { referalLinkService } from '@/services/referals'

const referrals = ref<ReferalLink[]>([])
const isLoading = ref(false)

async function loadReferrals() {
  isLoading.value = true
  try {
    const result = await referalLinkService.search(5, 0, { status: 'active' })
    referrals.value = result.items
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

onMounted(() => loadReferrals())
</script>

<template>
  <aside class="space-y-4">
    <h2 class="text-lg font-semibold">
      Актуальные рефералки
    </h2>

    <div
      v-if="isLoading"
      class="flex justify-center py-8"
    >
      <Loader2 class="h-6 w-6 animate-spin text-muted-foreground" />
    </div>

    <div
      v-else-if="referrals.length === 0"
      class="text-sm text-muted-foreground py-4"
    >
      Нет активных рефералок
    </div>

    <div
      v-else
      class="space-y-3"
    >
      <ReferralCard
        v-for="referral in referrals"
        :key="referral.id"
        :referral="referral"
      />
    </div>

    <router-link
      to="/referals"
      class="block text-center text-sm text-primary hover:underline"
    >
      Показать все →
    </router-link>
  </aside>
</template>
