<script setup lang="ts">
import type { RaffleCreateRequest } from '@/services/raffleService'
import { Typography } from 'itx-ui-kit'
import { Plus, Trash2 } from 'lucide-vue-next'
import { onMounted, ref } from 'vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { raffleService } from '@/services/raffleService'

const showModal = ref(false)
const confirmDeleteId = ref<number | null>(null)

const form = ref<RaffleCreateRequest>({
  title: '',
  description: '',
  prize: '',
  ticketCost: 10,
  maxTickets: 0,
  endsAt: '',
})

function resetForm() {
  form.value = {
    title: '',
    description: '',
    prize: '',
    ticketCost: 10,
    maxTickets: 0,
    endsAt: '',
  }
}

function openCreate() {
  resetForm()
  showModal.value = true
}

async function handleSubmit() {
  const data = {
    ...form.value,
    endsAt: new Date(form.value.endsAt).toISOString(),
  }
  const ok = await raffleService.create(data)
  if (ok) {
    showModal.value = false
    resetForm()
  }
}

async function handleDelete() {
  if (!confirmDeleteId.value)
    return
  await raffleService.delete(confirmDeleteId.value)
  confirmDeleteId.value = null
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

onMounted(() => {
  raffleService.getAll()
})
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex items-center justify-between">
        <Typography
          variant="h2"
          as="h1"
        >
          Розыгрыши
        </Typography>
        <Button
          size="sm"
          @click="openCreate"
        >
          <Plus class="h-4 w-4 mr-1" />
          Создать розыгрыш
        </Button>
      </div>

      <Card>
        <CardContent class="p-0">
          <div class="overflow-x-auto">
            <table class="w-full text-sm">
              <thead>
                <tr class="border-b">
                  <th class="text-left py-3 px-4 font-medium">
                    Название
                  </th>
                  <th class="text-left py-3 px-4 font-medium">
                    Приз
                  </th>
                  <th class="text-center py-3 px-4 font-medium">
                    Билет
                  </th>
                  <th class="text-center py-3 px-4 font-medium">
                    Лимит
                  </th>
                  <th class="text-left py-3 px-4 font-medium">
                    Завершение
                  </th>
                  <th class="text-center py-3 px-4 font-medium">
                    Статус
                  </th>
                  <th class="text-right py-3 px-4 font-medium">
                    Действия
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="raffle in raffleService.items.value"
                  :key="raffle.id"
                  class="border-b last:border-0 hover:bg-muted/50"
                >
                  <td class="py-3 px-4">
                    <div class="font-medium">
                      {{ raffle.title }}
                    </div>
                    <div
                      v-if="raffle.description"
                      class="text-xs text-muted-foreground mt-0.5"
                    >
                      {{ raffle.description }}
                    </div>
                  </td>
                  <td class="py-3 px-4">
                    {{ raffle.prize }}
                  </td>
                  <td class="py-3 px-4 text-center">
                    {{ raffle.ticketCost }} б.
                  </td>
                  <td class="py-3 px-4 text-center text-muted-foreground">
                    {{ raffle.maxTickets || '∞' }}
                  </td>
                  <td class="py-3 px-4 text-muted-foreground">
                    {{ formatDate(raffle.endsAt) }}
                  </td>
                  <td class="py-3 px-4 text-center">
                    <span
                      class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                      :class="raffle.status === 'ACTIVE'
                        ? 'bg-green-500/10 text-green-500'
                        : 'bg-muted text-muted-foreground'"
                    >
                      {{ raffle.status === 'ACTIVE' ? 'Активный' : 'Завершён' }}
                    </span>
                  </td>
                  <td class="py-3 px-4 text-right">
                    <Button
                      variant="ghost"
                      size="sm"
                      @click="confirmDeleteId = raffle.id"
                    >
                      <Trash2 class="h-4 w-4 text-destructive" />
                    </Button>
                  </td>
                </tr>
                <tr v-if="raffleService.items.value.length === 0 && !raffleService.isLoading.value">
                  <td
                    colspan="7"
                    class="py-8 text-center text-muted-foreground"
                  >
                    Розыгрышей пока нет
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>

      <!-- Confirm delete dialog -->
      <Teleport to="body">
        <div
          v-if="confirmDeleteId"
          class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
          @mousedown.self="confirmDeleteId = null"
        >
          <Card class="w-full max-w-sm">
            <CardContent class="p-6 space-y-4">
              <p class="text-sm">
                Удалить розыгрыш? Это действие нельзя отменить.
              </p>
              <div class="flex justify-end gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  @click="confirmDeleteId = null"
                >
                  Отмена
                </Button>
                <Button
                  variant="destructive"
                  size="sm"
                  @click="handleDelete"
                >
                  Удалить
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </Teleport>

      <!-- Create modal -->
      <Teleport to="body">
        <div
          v-if="showModal"
          class="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4"
          @mousedown.self="showModal = false"
        >
          <Card class="w-full max-w-lg">
            <CardHeader>
              <CardTitle>Новый розыгрыш</CardTitle>
            </CardHeader>
            <CardContent>
              <form
                class="space-y-4"
                @submit.prevent="handleSubmit"
              >
                <div>
                  <label class="text-sm font-medium mb-1 block">Название</label>
                  <input
                    v-model="form.title"
                    type="text"
                    required
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    placeholder="Розыгрыш подписки"
                  >
                </div>

                <div>
                  <label class="text-sm font-medium mb-1 block">Описание</label>
                  <textarea
                    v-model="form.description"
                    rows="2"
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background resize-none"
                    placeholder="Описание розыгрыша (необязательно)"
                  />
                </div>

                <div>
                  <label class="text-sm font-medium mb-1 block">Приз</label>
                  <input
                    v-model="form.prize"
                    type="text"
                    required
                    class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    placeholder="Месяц подписки"
                  >
                </div>

                <div class="grid grid-cols-3 gap-4">
                  <div>
                    <label class="text-sm font-medium mb-1 block">Стоимость билета</label>
                    <input
                      v-model.number="form.ticketCost"
                      type="number"
                      min="1"
                      required
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                  </div>
                  <div>
                    <label class="text-sm font-medium mb-1 block">Лимит билетов</label>
                    <input
                      v-model.number="form.maxTickets"
                      type="number"
                      min="0"
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                      placeholder="0 = без лимита"
                    >
                  </div>
                  <div>
                    <label class="text-sm font-medium mb-1 block">Завершение</label>
                    <input
                      v-model="form.endsAt"
                      type="datetime-local"
                      required
                      class="w-full border rounded-md px-3 py-2 text-sm bg-background"
                    >
                  </div>
                </div>

                <div class="flex justify-end gap-2 pt-2">
                  <Button
                    type="button"
                    variant="outline"
                    @click="showModal = false"
                  >
                    Отмена
                  </Button>
                  <Button type="submit">
                    Создать
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
      </Teleport>
    </div>
  </AdminLayout>
</template>
