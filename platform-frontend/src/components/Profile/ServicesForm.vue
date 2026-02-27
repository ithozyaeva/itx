<script setup lang="ts">
import type { Mentor, Service } from '@/models/profile'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Textarea } from '@/components/ui/textarea'
import { useUser } from '@/composables/useUser'
import { handleError } from '@/services/errorService'
import { profileService } from '@/services/profile'
import { Typography } from 'itx-ui-kit'
import { Edit, Loader2, Plus, Trash2 } from 'lucide-vue-next'
import { ref, watchEffect } from 'vue'

const services = ref<Service[]>([])

watchEffect(() => {
  services.value = [...(useUser<Mentor>().value?.services ?? [])]
})

const isEdit = ref<boolean>(false)
const isLoading = ref(false)

function addService() {
  services.value.push({
    name: '',
    price: 0,
    id: 0,
  })
}

// Удаление услуги
function removeService(index: number) {
  services.value.splice(index, 1)
}

async function handleSubmit() {
  isLoading.value = true
  try {
    await profileService.updateServices(services.value)
    isEdit.value = false
  }
  catch (error) {
    handleError(error)
  }
  finally {
    isLoading.value = false
  }
}

function toggleEdit() {
  services.value = [...useUser<Mentor>().value?.services || []]
  isEdit.value = !isEdit.value
}
</script>

<template>
  <div class="p-6 md:p-8 bg-card backdrop-blur-lg border border-border shadow-lg rounded-3xl h-fit">
    <div class="flex relative flex-col items-center space-y-3">
      <Edit class="absolute right-0 top-0 cursor-pointer text-muted-foreground hover:text-foreground" @click="toggleEdit" />
      <Typography variant="h3" as="h5">
        Услуги
      </Typography>
      <Table v-if="services.length > 0">
        <TableHeader>
          <TableRow>
            <TableHead>Описание</TableHead>
            <TableHead class="w-[150px]">
              Цена
            </TableHead>
            <TableHead class="w-[5px]">
              Действия
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="(service, index) in services" :key="index">
            <TableCell>
              <Textarea v-model="service.name" placeholder="Описание" :readonly="!isEdit" />
            </TableCell>
            <TableCell>
              <Input v-model="service.price" type="number" placeholder="Цена" :readonly="!isEdit" />
            </TableCell>
            <TableCell>
              <Button
                variant="ghost" size="sm" class="text-destructive" type="button"
                :disabled="!isEdit"
                @click="removeService(index)"
              >
                <Trash2 class="h-4 w-4" />
              </Button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <Button v-if="isEdit" variant="outline" size="sm" type="button" @click="addService">
        <Plus class="h-4 w-4 mr-1" /> Добавить услугу
      </Button>
      <div v-if="services.length === 0 && !isEdit">
        Нет услуг :(
      </div>
      <Button v-if="isEdit" class="mt-1 px-4 py-2 cursor-pointer transition duration-300 gap-2" :disabled="isLoading" @click="handleSubmit">
        <Loader2 v-if="isLoading" class="h-4 w-4 animate-spin" />
        Сохранить изменения
      </Button>
    </div>
  </div>
</template>
