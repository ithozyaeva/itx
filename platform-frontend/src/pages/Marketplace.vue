<script setup lang="ts">
import type { MarketplaceItem, MarketplaceItemStatus } from '@/models/marketplace'
import { Typography } from 'itx-ui-kit'
import {
  Loader2,
  Package,
  Plus,
  Trash2,
  User,
} from 'lucide-vue-next'
import { computed, onMounted, ref } from 'vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ErrorState from '@/components/common/ErrorState.vue'
import FormField from '@/components/common/FormField.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import {
  Dialog,
  DialogFooter,
  DialogHeader,
  DialogScrollContent,
  DialogTitle,
} from '@/components/ui/dialog'
import { required, useFormValidation } from '@/composables/useFormValidation'
import { isUserAdmin, useUser } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { marketplaceService } from '@/services/marketplace'

const items = ref<MarketplaceItem[]>([])
const total = ref(0)
const isLoading = ref(true)
const loadError = ref<string | null>(null)
const isSubmitting = ref(false)
const showCreateDialog = ref(false)
const activeStatus = ref<MarketplaceItemStatus | 'all'>('ACTIVE')

const user = useUser()
const isAdmin = isUserAdmin()

const { errors, validateAll, clearErrors } = useFormValidation({ title: [required('Введите название товара')] })

// Form state
const newTitle = ref('')
const newDescription = ref('')
const newPrice = ref('')
const newCity = ref('')
const newCanShip = ref(false)
const newCondition = ref<'NEW' | 'USED'>('NEW')
const newDefects = ref('')
const newPackageContents = ref('')
const newImage = ref<File | null>(null)

const statusTabs: { key: MarketplaceItemStatus | 'all', label: string }[] = [
  { key: 'all', label: 'Все' },
  { key: 'ACTIVE', label: 'Активные' },
  { key: 'RESERVED', label: 'Забронированные' },
  { key: 'SOLD', label: 'Проданные' },
]

const statusConfig: Record<MarketplaceItemStatus, { label: string, class: string }> = {
  ACTIVE: { label: 'Активно', class: 'bg-blue-500/10 text-blue-500' },
  RESERVED: { label: 'Забронировано', class: 'bg-yellow-500/10 text-yellow-500' },
  SOLD: { label: 'Продано', class: 'bg-green-500/10 text-green-500' },
  ARCHIVED: { label: 'В архиве', class: 'bg-muted text-muted-foreground' },
}

const filteredItems = computed(() => {
  if (activeStatus.value === 'all')
    return items.value
  return items.value.filter(i => i.status === activeStatus.value)
})

async function fetchItems() {
  isLoading.value = true
  loadError.value = null
  try {
    const res = await marketplaceService.getAll({ limit: 100 })
    items.value = res.items ?? []
    total.value = res.total
  }
  catch (error) {
    loadError.value = (await handleError(error)).message
  }
  finally {
    isLoading.value = false
  }
}

async function createItem() {
  if (!validateAll({ title: newTitle.value }))
    return
  isSubmitting.value = true
  try {
    await marketplaceService.create(
      {
        title: newTitle.value.trim(),
        description: newDescription.value.trim(),
        price: newPrice.value.trim(),
        city: newCity.value.trim(),
        canShip: newCanShip.value,
        condition: newCondition.value,
        defects: newDefects.value.trim(),
        packageContents: newPackageContents.value.trim(),
        contactTelegram: user.value?.tg ?? '',
      },
      newImage.value ?? undefined,
    )
    showCreateDialog.value = false
    resetForm()
    await fetchItems()
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isSubmitting.value = false
  }
}

function resetForm() {
  newTitle.value = ''
  newDescription.value = ''
  newPrice.value = ''
  newCity.value = ''
  newCanShip.value = false
  newCondition.value = 'NEW'
  newDefects.value = ''
  newPackageContents.value = ''
  newImage.value = null
  clearErrors()
}

async function requestPurchase(id: number) {
  try {
    await marketplaceService.requestPurchase(id)
    await fetchItems()
  }
  catch (error) {
    handleError(error)
  }
}

async function cancelPurchase(id: number) {
  try {
    await marketplaceService.cancelPurchase(id)
    await fetchItems()
  }
  catch (error) {
    handleError(error)
  }
}

async function markSold(id: number) {
  try {
    await marketplaceService.markSold(id)
    await fetchItems()
  }
  catch (error) {
    handleError(error)
  }
}

async function deleteItem(id: number) {
  try {
    await marketplaceService.remove(id)
    await fetchItems()
  }
  catch (error) {
    handleError(error)
  }
}

function isSeller(item: MarketplaceItem) {
  return user.value?.id === item.sellerId
}

function isBuyer(item: MarketplaceItem) {
  return user.value?.id === item.buyerId
}

function displayName(member: { firstName: string, lastName: string, tg: string }) {
  const name = [member.firstName, member.lastName].filter(Boolean).join(' ')
  return name || `@${member.tg}`
}

function onImageChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0] ?? null
  if (file) {
    const maxSize = 5 * 1024 * 1024 // 5MB
    if (file.size > maxSize) {
      handleError(new Error('Размер файла не должен превышать 5 МБ'))
      target.value = ''
      return
    }
  }
  newImage.value = file
}

onMounted(() => {
  fetchItems()
})
</script>

<template>
  <div class="container mx-auto px-4 py-6 md:py-8">
    <div class="flex items-center justify-between mb-6">
      <Typography
        variant="h2"
        as="h1"
      >
        Барахолка
      </Typography>
      <button
        class="flex items-center gap-2 px-4 py-2 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors"
        @click="showCreateDialog = true"
      >
        <Plus class="h-4 w-4" />
        Новое объявление
      </button>
    </div>

    <div
      v-if="isLoading"
      class="flex justify-center py-12"
    >
      <Loader2 class="h-8 w-8 animate-spin text-muted-foreground" />
    </div>

    <ErrorState
      v-else-if="loadError"
      :message="loadError"
      @retry="fetchItems"
    />

    <template v-else>
      <div class="flex gap-2 mb-6 flex-wrap">
        <button
          v-for="tab in statusTabs"
          :key="tab.key"
          class="px-3 py-1.5 rounded-full text-sm font-medium transition-colors"
          :class="activeStatus === tab.key
            ? 'bg-primary text-primary-foreground'
            : 'bg-card border border-border text-muted-foreground hover:text-foreground'"
          @click="activeStatus = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>

      <EmptyState
        v-if="filteredItems.length === 0"
        :icon="Package"
        title="Объявлений пока нет"
        description="Продайте или обменяйте ненужные вещи с участниками сообщества"
        action-label="Новое объявление"
        @action="showCreateDialog = true"
      />

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="item in filteredItems"
          :key="item.id"
          class="rounded-2xl border bg-card border-border overflow-hidden"
        >
          <!-- Image -->
          <img
            v-if="item.imagePath"
            :src="item.imagePath"
            :alt="item.title"
            loading="lazy"
            class="w-full h-48 object-cover rounded-t-2xl bg-muted"
          >
          <div
            v-else
            class="w-full h-48 bg-muted flex items-center justify-center rounded-t-2xl"
          >
            <Package class="h-12 w-12 text-muted-foreground opacity-40" aria-hidden="true" />
          </div>

          <div class="p-4">
            <!-- Title + status -->
            <div class="flex items-start justify-between gap-2 mb-2">
              <h3 class="font-medium text-sm leading-tight">
                {{ item.title }}
              </h3>
              <span
                class="shrink-0 px-2 py-0.5 rounded-full text-xs font-medium"
                :class="statusConfig[item.status].class"
              >
                {{ statusConfig[item.status].label }}
              </span>
            </div>

            <!-- Price -->
            <p class="font-bold text-sm mb-2">
              {{ item.price || 'Договорная' }}
            </p>

            <!-- City + shipping -->
            <div
              v-if="item.city || item.canShip"
              class="flex items-center gap-2 mb-2 flex-wrap"
            >
              <span
                v-if="item.city"
                class="text-xs text-muted-foreground"
              >{{ item.city }}</span>
              <span
                v-if="item.canShip"
                class="px-2 py-0.5 rounded-full text-xs bg-blue-500/10 text-blue-500"
              >Отправка</span>
            </div>

            <!-- Condition -->
            <div class="mb-2">
              <span
                class="px-2 py-0.5 rounded-full text-xs font-medium"
                :class="item.condition === 'NEW'
                  ? 'bg-green-500/10 text-green-500'
                  : 'bg-yellow-500/10 text-yellow-500'"
              >
                {{ item.condition === 'NEW' ? 'Новый' : 'Б/у' }}
              </span>
            </div>

            <!-- Defects -->
            <p
              v-if="item.defects"
              class="text-xs text-muted-foreground mb-1"
            >
              Дефекты: {{ item.defects }}
            </p>

            <!-- Package contents -->
            <p
              v-if="item.packageContents"
              class="text-xs text-muted-foreground mb-2"
            >
              Комплектация: {{ item.packageContents }}
            </p>

            <!-- Seller -->
            <div class="flex items-center gap-1.5 text-xs text-muted-foreground mb-1">
              <User class="h-3.5 w-3.5" />
              <span>Продавец: {{ displayName(item.seller) }}</span>
            </div>

            <!-- Buyer -->
            <div
              v-if="item.buyer"
              class="flex items-center gap-1.5 text-xs text-muted-foreground mb-1"
            >
              <User class="h-3.5 w-3.5" />
              <span>Покупатель: {{ displayName(item.buyer) }}</span>
            </div>

            <!-- Contact -->
            <div
              v-if="item.contactTelegram"
              class="mt-2 mb-3"
            >
              <a
                :href="`https://t.me/${item.contactTelegram.replace('@', '')}`"
                target="_blank"
                class="text-xs text-accent hover:underline"
              >
                @{{ item.contactTelegram.replace('@', '') }}
              </a>
            </div>

            <!-- Actions -->
            <div class="flex flex-wrap gap-2 mt-3">
              <!-- ACTIVE: buy (not seller) -->
              <button
                v-if="item.status === 'ACTIVE' && !isSeller(item)"
                class="px-3 py-1.5 rounded-lg text-xs font-medium bg-primary text-primary-foreground hover:bg-primary/90 transition-colors"
                @click="requestPurchase(item.id)"
              >
                Хочу купить
              </button>

              <!-- RESERVED: cancel (buyer) -->
              <button
                v-if="item.status === 'RESERVED' && isBuyer(item)"
                class="px-3 py-1.5 rounded-lg text-xs font-medium bg-muted text-muted-foreground hover:text-foreground transition-colors"
                @click="cancelPurchase(item.id)"
              >
                Отменить бронь
              </button>

              <!-- RESERVED: cancel reservation (seller) -->
              <ConfirmDialog
                v-if="item.status === 'RESERVED' && isSeller(item)"
                title="Снять бронь?"
                :description="`Бронь покупателя ${item.buyer ? displayName(item.buyer) : ''} будет отменена, и объявление вернётся в активные.`"
                confirm-label="Снять бронь"
                @confirm="cancelPurchase(item.id)"
              >
                <template #trigger>
                  <button
                    class="px-3 py-1.5 rounded-lg text-xs font-medium bg-muted text-muted-foreground hover:text-foreground transition-colors"
                  >
                    Снять бронь
                  </button>
                </template>
              </ConfirmDialog>

              <!-- RESERVED: confirm sale (seller) -->
              <button
                v-if="item.status === 'RESERVED' && isSeller(item)"
                class="px-3 py-1.5 rounded-lg text-xs font-medium bg-green-600 text-white hover:bg-green-700 transition-colors"
                @click="markSold(item.id)"
              >
                Подтвердить продажу
              </button>

              <!-- Delete (ACTIVE seller, or admin) -->
              <ConfirmDialog
                v-if="(item.status === 'ACTIVE' && isSeller(item)) || isAdmin"
                title="Удалить объявление?"
                description="Объявление будет удалено безвозвратно."
                confirm-label="Удалить"
                @confirm="deleteItem(item.id)"
              >
                <template #trigger>
                  <button
                    class="flex items-center gap-1 px-3 py-1.5 rounded-lg text-xs font-medium text-red-500 hover:bg-red-500/10 transition-colors ml-auto"
                    aria-label="Удалить объявление"
                  >
                    <Trash2 class="h-3.5 w-3.5" />
                  </button>
                </template>
              </ConfirmDialog>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Create dialog -->
    <Dialog v-model:open="showCreateDialog">
      <DialogScrollContent>
        <DialogHeader>
          <DialogTitle>Новое объявление</DialogTitle>
        </DialogHeader>

        <form
          class="space-y-4"
          @submit.prevent="createItem"
        >
          <FormField
            label="Название"
            :error="errors.title"
            required
            html-for="create-item-title"
          >
            <input
              id="create-item-title"
              v-model="newTitle"
              type="text"
              required
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="Название товара"
            >
          </FormField>

          <div>
            <label class="block text-sm font-medium mb-1">Описание</label>
            <textarea
              v-model="newDescription"
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary min-h-20 resize-none"
              placeholder="Подробное описание..."
            />
          </div>

          <div>
            <label class="block text-sm font-medium mb-1">Цена</label>
            <input
              v-model="newPrice"
              type="text"
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="Договорная"
            >
          </div>

          <div>
            <label class="block text-sm font-medium mb-1">Город</label>
            <input
              v-model="newCity"
              type="text"
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="Москва"
            >
          </div>

          <div class="flex items-center gap-2">
            <input
              id="canShip"
              v-model="newCanShip"
              type="checkbox"
              class="h-4 w-4 rounded border-border"
            >
            <label
              for="canShip"
              class="text-sm"
            >Возможна отправка почтой</label>
          </div>

          <div>
            <label class="block text-sm font-medium mb-2">Состояние</label>
            <div class="flex gap-2">
              <button
                type="button"
                class="px-3 py-1.5 rounded-xl text-sm font-medium border transition-colors"
                :class="newCondition === 'NEW'
                  ? 'bg-primary text-primary-foreground border-primary'
                  : 'bg-card border-border text-muted-foreground hover:text-foreground'"
                @click="newCondition = 'NEW'"
              >
                Новый
              </button>
              <button
                type="button"
                class="px-3 py-1.5 rounded-xl text-sm font-medium border transition-colors"
                :class="newCondition === 'USED'
                  ? 'bg-primary text-primary-foreground border-primary'
                  : 'bg-card border-border text-muted-foreground hover:text-foreground'"
                @click="newCondition = 'USED'"
              >
                Б/у
              </button>
            </div>
          </div>

          <div v-if="newCondition === 'USED'">
            <label class="block text-sm font-medium mb-1">Дефекты</label>
            <textarea
              v-model="newDefects"
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary min-h-16 resize-none"
              placeholder="Опишите дефекты..."
            />
          </div>

          <div>
            <label class="block text-sm font-medium mb-1">Комплектация</label>
            <textarea
              v-model="newPackageContents"
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary min-h-16 resize-none"
              placeholder="Что входит в комплект..."
            />
          </div>

          <div>
            <label class="block text-sm font-medium mb-1">Фото</label>
            <input
              type="file"
              accept="image/*"
              class="w-full text-sm text-muted-foreground file:mr-3 file:py-1 file:px-3 file:rounded-lg file:border-0 file:text-sm file:font-medium file:bg-primary file:text-primary-foreground hover:file:bg-primary/90"
              @change="onImageChange"
            >
          </div>

          <DialogFooter>
            <button
              type="submit"
              class="px-4 py-2 rounded-xl bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors disabled:opacity-50"
              :disabled="!newTitle.trim() || isSubmitting"
            >
              <Loader2
                v-if="isSubmitting"
                class="h-4 w-4 animate-spin inline mr-1"
              />
              Опубликовать
            </button>
          </DialogFooter>
        </form>
      </DialogScrollContent>
    </Dialog>
  </div>
</template>
