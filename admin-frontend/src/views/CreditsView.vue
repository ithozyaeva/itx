<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import Plus from '~icons/lucide/plus'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import CreditsAwardModal from '@/components/modals/CreditsAwardModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Typography } from '@/components/ui/typography'
import { creditsService } from '@/services/creditsService'

const isAwardModalOpen = ref(false)
const usernameFilter = ref('')

const reasonLabels: Record<string, string> = {
  referal_conversion: 'Конверсия по реф-ссылке',
  referral_purchase_first: 'Реферал впервые оформил подписку',
  referral_purchase_recurring: 'Реферал — ежемесячная активность',
  subscription_purchase: 'Покупка подписки',
  admin_manual: 'Начислено админом',
}

function applyMemberFilter() {
  creditsService.applyFilters({ username: usernameFilter.value || undefined })
}

function resetFilters() {
  usernameFilter.value = ''
  creditsService.applyFilters({})
}

onMounted(creditsService.search)
onUnmounted(creditsService.clearPagination)
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <Typography variant="h2" as="h1">
          Реферальные кредиты
        </Typography>
        <Button
          v-permission="'can_edit_admin_points'"
          @click="isAwardModalOpen = true"
        >
          <Plus class="mr-2 h-4 w-4" />
          Начислить кредиты
        </Button>
      </div>

      <Card class="p-4 rounded-lg">
        <div class="flex items-center gap-3">
          <Input
            v-model="usernameFilter"
            placeholder="TG username"
            class="w-[200px]"
            @keydown.enter="applyMemberFilter"
          />
          <Button
            size="sm"
            @click="applyMemberFilter"
          >
            Применить
          </Button>
          <Button
            variant="outline"
            size="sm"
            @click="resetFilters"
          >
            Сбросить
          </Button>
        </div>
      </Card>

      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Участник</TableHead>
                <TableHead>Сумма</TableHead>
                <TableHead>Причина</TableHead>
                <TableHead>Описание</TableHead>
                <TableHead>Дата</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow
                v-if="creditsService.items.value.total === 0"
                class="h-24"
              >
                <TableCell
                  colspan="5"
                  class="text-center"
                >
                  Транзакции не найдены
                </TableCell>
              </TableRow>
              <TableRow
                v-for="tx in creditsService.items.value.items"
                :key="tx.id"
              >
                <TableCell>
                  {{ tx.memberFirstName ?? '' }} {{ tx.memberLastName ?? '' }}
                  <span class="text-muted-foreground">@{{ tx.memberUsername }}</span>
                </TableCell>
                <TableCell
                  class="font-medium"
                  :class="tx.amount > 0 ? 'text-green-600' : 'text-red-600'"
                >
                  {{ tx.amount > 0 ? '+' : '' }}{{ tx.amount }}
                </TableCell>
                <TableCell>{{ reasonLabels[tx.reason] ?? tx.reason }}</TableCell>
                <TableCell class="max-w-[240px] truncate">
                  {{ tx.description }}
                </TableCell>
                <TableCell class="whitespace-nowrap">
                  {{ new Date(tx.createdAt).toLocaleString() }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <div class="mt-4 flex justify-end">
        <Pagination
          v-slot="{ page }"
          :items-per-page="20"
          :total="creditsService.items.value.total"
          :sibling-count="1"
          show-edges
          :default-page="1"
        >
          <PaginationList
            v-slot="{ items }"
            class="flex items-center gap-1"
          >
            <PaginationFirst />
            <PaginationPrev />

            <template v-for="(item, index) in items">
              <PaginationListItem
                v-if="item.type === 'page'"
                :key="index"
                :value="item.value"
                as-child
              >
                <Button
                  class="w-10 h-10 p-0"
                  :variant="item.value === page ? 'default' : 'outline'"
                  @click="creditsService.changePagination(item.value)"
                >
                  {{ item.value }}
                </Button>
              </PaginationListItem>
              <PaginationEllipsis
                v-else
                :key="item.type"
                :index="index"
              />
            </template>

            <PaginationNext />
            <PaginationLast />
          </PaginationList>
        </Pagination>
      </div>
    </div>

    <CreditsAwardModal
      v-model:is-open="isAwardModalOpen"
      @saved="creditsService.search"
    />
  </AdminLayout>
</template>
