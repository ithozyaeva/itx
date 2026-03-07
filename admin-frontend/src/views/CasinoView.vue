<script setup lang="ts">
import { Typography } from 'itx-ui-kit'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { casinoService } from '@/services/casinoService'

const usernameFilter = ref('')
const gameFilter = ref<string>('all')

const gameLabels: Record<string, string> = {
  coin_flip: 'Монетка',
  dice_roll: 'Кости',
  wheel: 'Колесо',
}

const stats = computed(() => casinoService.stats.value)

function applyFilters() {
  casinoService.applyFilters({
    username: usernameFilter.value || undefined,
    game: gameFilter.value === 'all' ? undefined : gameFilter.value,
  })
}

function resetFilters() {
  usernameFilter.value = ''
  gameFilter.value = 'all'
  casinoService.applyFilters({})
}

function formatNumber(n: number) {
  return n.toLocaleString('ru-RU')
}

onMounted(() => {
  casinoService.getStats()
  casinoService.searchBets()
})

onUnmounted(casinoService.clearPagination)
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Казино
      </Typography>

      <div
        v-if="stats"
        class="grid grid-cols-2 gap-4 md:grid-cols-5"
      >
        <Card>
          <CardContent class="p-4">
            <div class="text-sm text-muted-foreground">
              Всего ставок
            </div>
            <div class="text-2xl font-bold">
              {{ formatNumber(stats.totalBets) }}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent class="p-4">
            <div class="text-sm text-muted-foreground">
              Оборот
            </div>
            <div class="text-2xl font-bold">
              {{ formatNumber(stats.totalWagered) }}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent class="p-4">
            <div class="text-sm text-muted-foreground">
              Выплаты
            </div>
            <div class="text-2xl font-bold">
              {{ formatNumber(stats.totalPayout) }}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent class="p-4">
            <div class="text-sm text-muted-foreground">
              Прибыль
            </div>
            <div class="text-2xl font-bold" :class="stats.houseProfit >= 0 ? 'text-green-600' : 'text-red-600'">
              {{ formatNumber(stats.houseProfit) }}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardContent class="p-4">
            <div class="text-sm text-muted-foreground">
              Игроков
            </div>
            <div class="text-2xl font-bold">
              {{ formatNumber(stats.uniquePlayers) }}
            </div>
          </CardContent>
        </Card>
      </div>

      <Card v-if="stats && stats.gameStats.length > 0">
        <CardContent class="p-4">
          <div class="mb-3 text-sm font-medium text-muted-foreground">
            Статистика по играм
          </div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Игра</TableHead>
                <TableHead>Ставок</TableHead>
                <TableHead>Оборот</TableHead>
                <TableHead>Выплаты</TableHead>
                <TableHead>Прибыль</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow
                v-for="gs in stats.gameStats"
                :key="gs.game"
              >
                <TableCell>{{ gameLabels[gs.game] ?? gs.game }}</TableCell>
                <TableCell>{{ formatNumber(gs.totalBets) }}</TableCell>
                <TableCell>{{ formatNumber(gs.totalWagered) }}</TableCell>
                <TableCell>{{ formatNumber(gs.totalPayout) }}</TableCell>
                <TableCell :class="gs.houseProfit >= 0 ? 'text-green-600' : 'text-red-600'">
                  {{ formatNumber(gs.houseProfit) }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card class="p-4 rounded-lg">
        <div class="flex items-center gap-3">
          <Input
            v-model="usernameFilter"
            placeholder="TG username"
            class="w-[200px]"
            @keydown.enter="applyFilters"
          />
          <Select v-model="gameFilter">
            <SelectTrigger class="w-[180px]">
              <SelectValue placeholder="Все игры" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">
                Все игры
              </SelectItem>
              <SelectItem value="coin_flip">
                Монетка
              </SelectItem>
              <SelectItem value="dice_roll">
                Кости
              </SelectItem>
              <SelectItem value="wheel">
                Колесо
              </SelectItem>
            </SelectContent>
          </Select>
          <Button
            size="sm"
            @click="applyFilters"
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
                <TableHead>Игрок</TableHead>
                <TableHead>Игра</TableHead>
                <TableHead>Ставка</TableHead>
                <TableHead>Выбор</TableHead>
                <TableHead>Результат</TableHead>
                <TableHead>Множитель</TableHead>
                <TableHead>Выплата</TableHead>
                <TableHead>Профит</TableHead>
                <TableHead>Дата</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow
                v-if="casinoService.items.value.total === 0"
                class="h-24"
              >
                <TableCell
                  colspan="9"
                  class="text-center"
                >
                  Ставки не найдены
                </TableCell>
              </TableRow>
              <TableRow
                v-for="bet in casinoService.items.value.items"
                :key="bet.id"
              >
                <TableCell>
                  {{ bet.memberFirstName ?? '' }} {{ bet.memberLastName ?? '' }}
                  <span class="text-muted-foreground">@{{ bet.memberUsername }}</span>
                </TableCell>
                <TableCell>{{ gameLabels[bet.game] ?? bet.game }}</TableCell>
                <TableCell class="font-medium">
                  {{ bet.betAmount }}
                </TableCell>
                <TableCell>{{ bet.betChoice }}</TableCell>
                <TableCell>{{ bet.result }}</TableCell>
                <TableCell>x{{ bet.multiplier }}</TableCell>
                <TableCell>{{ bet.payout }}</TableCell>
                <TableCell :class="bet.profit >= 0 ? 'text-green-600' : 'text-red-600'">
                  {{ bet.profit }}
                </TableCell>
                <TableCell class="whitespace-nowrap">
                  {{ new Date(bet.createdAt).toLocaleString() }}
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
          :total="casinoService.items.value.total"
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
                  @click="casinoService.changePagination(item.value)"
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
  </AdminLayout>
</template>
