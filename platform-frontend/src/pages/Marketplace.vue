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
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { isUserAdmin, useUser } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { marketplaceService } from '@/services/marketplace'

const items = ref<MarketplaceItem[]>([])
const total = ref(0)
const isLoading = ref(true)
const isSubmitting = ref(false)
const showCreateDialog = ref(false)
const activeStatus = ref<MarketplaceItemStatus | 'all'>('all')

const user = useUser()
const isAdmin = isUserAdmin()

// Form state
const newTitle = ref('')
const newDescription = ref('')
const newPrice = ref('')
const newCity = ref('')
const newCanShip = ref(false)
const newCondition = ref<'NEW' | 'USED'>('NEW')
const newDefects = ref('')
const newPackageContents = ref('')
const newContactTelegram = ref('')
const newContactEmail = ref('')
const newContactPhone = ref('')
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
  try {
    const res = await marketplaceService.getAll({ limit: 100 })
    items.value = res.items ?? []
    total.value = res.total
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

async function createItem() {
  if (!newTitle.value.trim())
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
        contactTelegram: newContactTelegram.value.trim(),
        contactEmail: newContactEmail.value.trim(),
        contactPhone: newContactPhone.value.trim(),
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
  newContactTelegram.value = ''
  newContactEmail.value = ''
  newContactPhone.value = ''
  newImage.value = null
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
  newImage.value = target.files?.[0] ?? null
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

      <div
        v-if="filteredItems.length === 0"
        class="text-center py-12 text-muted-foreground"
      >
        <Package class="h-12 w-12 mx-auto mb-3 opacity-50" />
        <p>Объявлений пока нет</p>
      </div>

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
            class="w-full h-48 object-cover rounded-t-2xl"
          >
          <div
            v-else
            class="w-full h-48 bg-muted flex items-center justify-center rounded-t-2xl"
          >
            <Package class="h-12 w-12 text-muted-foreground opacity-40" />
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

            <!-- Contacts -->
            <div
              v-if="item.contactTelegram || item.contactEmail || item.contactPhone"
              class="mt-2 mb-3 space-y-0.5"
            >
              <p
                v-if="item.contactTelegram"
                class="text-xs text-muted-foreground"
              >
                TG: {{ item.contactTelegram }}
              </p>
              <p
                v-if="item.contactEmail"
                class="text-xs text-muted-foreground"
              >
                Email: {{ item.contactEmail }}
              </p>
              <p
                v-if="item.contactPhone"
                class="text-xs text-muted-foreground"
              >
                Тел: {{ item.contactPhone }}
              </p>
            </div>

            <!-- Actions -->
            <div class="flex flex-wrap gap-2 mt-3">
              <!-- ACTIVE: buy (not seller) -->
              <button
                v-if="item.status === 'ACTIVE' && !isSeller(item)"
                class="px-3 py-1 rounded-lg text-xs font-medium bg-primary text-primary-foreground hover:bg-primary/90 transition-colors"
                @click="requestPurchase(item.id)"
              >
                Хочу купить
              </button>

              <!-- RESERVED: cancel (buyer) -->
              <button
                v-if="item.status === 'RESERVED' && isBuyer(item)"
                class="px-3 py-1 rounded-lg text-xs font-medium bg-muted text-muted-foreground hover:text-foreground transition-colors"
                @click="cancelPurchase(item.id)"
              >
                Отменить бронь
              </button>

              <!-- RESERVED: confirm sale (seller) -->
              <button
                v-if="item.status === 'RESERVED' && isSeller(item)"
                class="px-3 py-1 rounded-lg text-xs font-medium bg-green-600 text-white hover:bg-green-700 transition-colors"
                @click="markSold(item.id)"
              >
                Подтвердить продажу
              </button>

              <!-- Delete (ACTIVE seller, or admin) -->
              <button
                v-if="(item.status === 'ACTIVE' && isSeller(item)) || isAdmin"
                class="flex items-center gap-1 px-3 py-1 rounded-lg text-xs font-medium text-red-500 hover:bg-red-500/10 transition-colors ml-auto"
                @click="deleteItem(item.id)"
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
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Новое объявление</DialogTitle>
        </DialogHeader>

        <form
          class="space-y-4"
          @submit.prevent="createItem"
        >
          <div>
            <label class="block text-sm font-medium mb-1">Название *</label>
            <input
              v-model="newTitle"
              type="text"
              required
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="Название товара"
            >
          </div>

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
            <label class="block text-sm font-medium mb-1">Telegram контакт</label>
            <input
              v-model="newContactTelegram"
              type="text"
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="@username"
            >
          </div>

          <div>
            <label class="block text-sm font-medium mb-1">Email</label>
            <input
              v-model="newContactEmail"
              type="email"
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="email@example.com"
            >
          </div>

          <div>
            <label class="block text-sm font-medium mb-1">Телефон</label>
            <input
              v-model="newContactPhone"
              type="tel"
              class="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="+7 900 000 00 00"
            >
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
      </DialogContent>
    </Dialog>
  </div>
</template>
