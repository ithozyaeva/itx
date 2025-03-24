<script setup lang="ts">
import type { ServiceFormData } from '@/models/mentors'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Textarea } from '@/components/ui/textarea'
import { Plus, Trash2 } from 'lucide-vue-next'

const props = defineProps<{
  services: ServiceFormData[]
}>()

const emit = defineEmits(['update:services'])

function addService() {
  const newServices = [...props.services]
  newServices.push({
    name: '',
    price: 0,
  })
  emit('update:services', newServices)
}

// Удаление услуги
function removeService(index: number) {
  const newServices = [...props.services]
  newServices.splice(index, 1)
  emit('update:services', newServices)
}
</script>

<template>
  <Card>
    <CardHeader class="flex flex-row items-center justify-between">
      <CardTitle>Услуги</CardTitle>
      <Button variant="outline" size="sm" type="button" @click="addService">
        <Plus class="h-4 w-4 mr-1" /> Добавить услугу
      </Button>
    </CardHeader>
    <CardContent>
      <Table v-if="services.length > 0">
        <TableHeader>
          <TableRow>
            <TableHead>Описание</TableHead>
            <TableHead class="w-[250px]">
              Цена
            </TableHead>
            <TableHead class="w-[100px]">
              Действия
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="(service, index) in services" :key="index">
            <TableCell>
              <Textarea v-model="service.name" placeholder="Описание" />
            </TableCell>
            <TableCell>
              <Input
                v-model="service.price"
                type="number"
                placeholder="Цена"
              />
            </TableCell>
            <TableCell>
              <Button
                variant="ghost"
                size="sm"
                class="text-destructive"
                type="button"
                @click="removeService(index)"
              >
                <Trash2 class="h-4 w-4" />
              </Button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <div v-else class="text-center py-4 text-gray-500">
        Нет добавленных услуг
      </div>
    </CardContent>
  </Card>
</template>
