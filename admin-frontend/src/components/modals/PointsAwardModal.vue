<script setup lang="ts">
import type { Member } from '@/models/members'
import type { Registry } from '@/models/registry'
import { ref, watch } from 'vue'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import api from '@/lib/api'
import { pointsService } from '@/services/pointsService'

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits(['update:isOpen', 'saved'])

const selectedMemberId = ref<number | null>(null)
const memberSearch = ref('')
const memberResults = ref<Member[]>([])
const selectedMemberName = ref('')
const amount = ref<number>(0)
const description = ref('')
const isSearching = ref(false)
const showDropdown = ref(false)

let searchTimeout: ReturnType<typeof setTimeout> | null = null

watch(() => props.isOpen, (open) => {
  if (open) {
    selectedMemberId.value = null
    memberSearch.value = ''
    memberResults.value = []
    selectedMemberName.value = ''
    amount.value = 0
    description.value = ''
  }
})

async function searchMembers(query: string) {
  if (query.length < 2) {
    memberResults.value = []
    showDropdown.value = false
    return
  }

  if (searchTimeout)
    clearTimeout(searchTimeout)

  searchTimeout = setTimeout(async () => {
    try {
      isSearching.value = true
      const response = await api.get('members', {
        searchParams: { username: query, limit: 10 },
      }).json<Registry<Member>>()
      memberResults.value = response.items
      showDropdown.value = true
    }
    finally {
      isSearching.value = false
    }
  }, 300)
}

function selectMember(member: Member) {
  selectedMemberId.value = member.id
  selectedMemberName.value = `${member.firstName ?? ''} ${member.lastName ?? ''} (@${member.tg})`.trim()
  memberSearch.value = selectedMemberName.value
  showDropdown.value = false
}

function handleClose() {
  emit('update:isOpen', false)
}

async function handleSubmit() {
  if (!selectedMemberId.value || amount.value <= 0)
    return

  const success = await pointsService.award({
    memberId: selectedMemberId.value,
    amount: amount.value,
    description: description.value,
  })

  if (success) {
    emit('saved')
    handleClose()
  }
}
</script>

<template>
  <Dialog
    :open="isOpen"
    @update:open="handleClose"
  >
    <DialogContent class="sm:max-w-[450px]">
      <DialogHeader>
        <DialogTitle>Начислить баллы</DialogTitle>
      </DialogHeader>

      <div class="space-y-4">
        <div class="relative">
          <Label>Участник</Label>
          <Input
            v-model="memberSearch"
            placeholder="Поиск по имени или username"
            class="mt-1"
            @input="searchMembers(memberSearch)"
            @focus="memberResults.length > 0 && (showDropdown = true)"
          />
          <div
            v-if="showDropdown && memberResults.length > 0"
            class="absolute z-50 mt-1 w-full bg-popover border rounded-md shadow-md max-h-48 overflow-auto"
          >
            <button
              v-for="member in memberResults"
              :key="member.id"
              class="w-full px-3 py-2 text-left text-sm hover:bg-accent cursor-pointer"
              @mousedown.prevent="selectMember(member)"
            >
              {{ member.firstName ?? '' }} {{ member.lastName ?? '' }}
              <span class="text-muted-foreground">@{{ member.tg }}</span>
            </button>
          </div>
        </div>

        <div>
          <Label>Сумма баллов</Label>
          <Input
            v-model.number="amount"
            type="number"
            min="1"
            placeholder="Количество баллов"
            class="mt-1"
          />
        </div>

        <div>
          <Label>Описание</Label>
          <Input
            v-model="description"
            placeholder="Причина начисления"
            class="mt-1"
          />
        </div>
      </div>

      <DialogFooter>
        <Button
          variant="outline"
          @click="handleClose"
        >
          Отмена
        </Button>
        <Button
          :disabled="!selectedMemberId || amount <= 0 || pointsService.isLoading.value"
          @click="handleSubmit"
        >
          Начислить
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
