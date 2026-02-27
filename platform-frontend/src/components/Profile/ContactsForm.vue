<script setup lang="ts">
import type { Contacts, Mentor } from '@/models/profile'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useUser } from '@/composables/useUser'
import { profileService } from '@/services/profile'
import { Typography } from 'itx-ui-kit'
import { Edit, Plus, Trash2 } from 'lucide-vue-next'
import { ref, watchEffect } from 'vue'

const isEdit = ref<boolean>(false)
const contacts = ref<Contacts[]>([])

watchEffect(() => {
  contacts.value = [...(useUser<Mentor>().value?.contacts ?? [])]
})
function addContact() {
  contacts.value.push({
    type: 1, // Тип по умолчанию, TODO: в базе поменять int на varchar
    link: '',
    id: 0,
  })
}

// Удаление контакта
function removeContact(index: number) {
  contacts.value.splice(index, 1)
}
function toggleEdit() {
  contacts.value = [...useUser<Mentor>().value?.contacts || []]
  isEdit.value = !isEdit.value
}
async function handleSubmit() {
  await profileService.updateContacts(contacts.value)
  isEdit.value = false
}
</script>

<template>
  <div class="p-6 md:p-8 bg-card backdrop-blur-lg border border-border shadow-lg rounded-3xl h-fit">
    <div class="flex relative flex-col items-center space-y-3">
      <Edit class="absolute right-0 top-0 cursor-pointer text-muted-foreground hover:text-foreground" @click="toggleEdit" />
      <Typography variant="h3" as="h5">
        Контакты
      </Typography>
      <Table v-if="contacts.length > 0">
        <TableHeader>
          <TableRow>
            <TableHead class="w-[250px]">
              Тип
            </TableHead>
            <TableHead>Значение</TableHead>
            <TableHead class="w-[100px]">
              Действия
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-for="(contact, index) in contacts" :key="index">
            <TableCell>
              <Select v-model="contact.type" :disabled="!isEdit">
                <SelectTrigger>
                  <SelectValue placeholder="Выберите тип" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem :value="1">
                    Telegram
                  </SelectItem>
                  <SelectItem :value="2">
                    Email
                  </SelectItem>
                  <SelectItem :value="3">
                    Телефон
                  </SelectItem>
                  <SelectItem :value="4">
                    Другое
                  </SelectItem>
                </SelectContent>
              </Select>
            </TableCell>
            <TableCell>
              <Input v-model="contact.link" placeholder="Введите ссылку (если это ТГ прям с https)" :readonly="!isEdit" />
            </TableCell>
            <TableCell>
              <Button :disabled="!isEdit" variant="ghost" size="sm" class="text-destructive" type="button" @click="removeContact(index)">
                <Trash2 class="h-4 w-4" />
              </Button>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <Button v-if="isEdit" variant="outline" size="sm" type="button" @click="addContact">
        <Plus class="h-4 w-4 mr-1" /> Добавить контакт
      </Button>
      <div v-if="contacts.length === 0 && !isEdit">
        Нет контактов :(
      </div>
      <Button
        v-if="isEdit"
        class="mt-1 px-4 py-2 cursor-pointer transition duration-300"
        @click="handleSubmit"
      >
        Сохранить изменения
      </Button>
    </div>
  </div>
</template>
