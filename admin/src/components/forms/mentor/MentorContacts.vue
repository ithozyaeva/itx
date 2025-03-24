<script setup lang="ts">
import type { ContactFormData } from '@/models/mentors'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Plus, Trash2 } from 'lucide-vue-next'

const props = defineProps<{
  contacts: ContactFormData[]
}>()

const emit = defineEmits(['update:contacts'])

// Добавление нового контакта
function addContact() {
  const newContacts = [...props.contacts]
  newContacts.push({
    type: 1, // Тип по умолчанию
    link: '',
  })
  emit('update:contacts', newContacts)
}

// Удаление контакта
function removeContact(index: number) {
  const newContacts = [...props.contacts]
  newContacts.splice(index, 1)
  emit('update:contacts', newContacts)
}
</script>

<template>
  <Card>
    <CardHeader class="flex flex-row items-center justify-between">
      <CardTitle>Контакты</CardTitle>
      <Button variant="outline" size="sm" type="button" @click="addContact">
        <Plus class="h-4 w-4 mr-1" /> Добавить контакт
      </Button>
    </CardHeader>
    <CardContent>
      <Table v-if="contacts.length > 0">
        <TableHeader>
          <TableRow>
            <TableHead class="w-[250px]">
              Тип
            </TableHead>
            <TableHead>Ссылка</TableHead>
            <TableHead class="w-[100px]">
              Действия
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="(contact, index) in contacts" :key="index">
            <TableCell>
              <Select v-model="contact.type">
                <SelectTrigger>
                  <SelectValue placeholder="Выберите тип" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem :value="1">
                    Telegram
                  </SelectItem>
                  <SelectItem :value="2">
                    WhatsApp
                  </SelectItem>
                  <SelectItem :value="3">
                    Email
                  </SelectItem>
                  <SelectItem :value="4">
                    Телефон
                  </SelectItem>
                  <SelectItem :value="5">
                    Другое
                  </SelectItem>
                </SelectContent>
              </Select>
            </TableCell>
            <TableCell>
              <Input v-model="contact.link" placeholder="Введите ссылку" />
            </TableCell>
            <TableCell>
              <Button
                variant="ghost"
                size="sm"
                class="text-destructive"
                type="button"
                @click="removeContact(index)"
              >
                <Trash2 class="h-4 w-4" />
              </Button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <div v-else class="text-center py-4 text-gray-500">
        Нет добавленных контактов
      </div>
    </CardContent>
  </Card>
</template>
