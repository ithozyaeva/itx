<script setup lang="ts">
import type { MemberRole } from '@/models/members'
import { reactive } from 'vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { useDictionary } from '@/composables/useDictionary'
import Card from './ui/card/Card.vue'

interface MemberSearchFilters {
  username: string
  roles: MemberRole[]
}

const emit = defineEmits<{
  apply: [filters: MemberSearchFilters]
}>()

const localFilters = reactive<MemberSearchFilters>({
  username: '',
  roles: [],
})

const { memberRoles } = useDictionary<MemberRole>(['memberRoles'])

function applyFilters() {
  emit('apply', { ...localFilters })
}

function resetFilters() {
  localFilters.username = ''
  localFilters.roles = []
  emit('apply', { ...localFilters })
}
</script>

<template>
  <Card class="p-4 rounded-lg mb-6">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <Label for="username" class="text-sm font-medium">TG username</Label>
        <Input
          id="username"
          v-model="localFilters.username"
          placeholder="Поиск по TG"
          class="mt-1"
        />
      </div>

      <div>
        <Label for="roles" class="text-sm font-medium">Роли</Label>
        <div class="mt-1">
          <Select id="roles" v-model="localFilters.roles" multiple>
            <SelectTrigger>
              <SelectValue placeholder="Выберите роли" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem
                v-for="role in memberRoles"
                :key="role.value"
                :value="role.value"
              >
                {{ role.label }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      <!-- actions -->
      <div class="md:col-span-2 flex justify-end gap-2">
        <Button variant="outline" size="sm" @click="resetFilters">
          Сбросить
        </Button>
        <Button size="sm" @click="applyFilters">
          Применить
        </Button>
      </div>
    </div>
  </Card>
</template>
