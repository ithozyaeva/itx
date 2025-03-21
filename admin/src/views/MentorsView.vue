<script setup lang="ts">
import type { Mentor } from '@/models/mentors'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import MentorModal from '@/components/modals/MentorModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { mentorService } from '@/services/mentorService'
import { onMounted, ref } from 'vue'

import Pencil from '~icons/lucide/pencil'
import Trash from '~icons/lucide/trash'

const isModalOpen = ref(false)
const selectedMentor = ref<Mentor | null>(null)

/**
 * Обработчик изменения ментора.
 *
 * @param mentor - Ментор.
 */
function handleEditMentor(mentor: Mentor) {
  selectedMentor.value = mentor
  isModalOpen.value = true
}

onMounted(mentorService.search)
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <h1 class="text-3xl font-bold">
          Менторы
        </h1>
      </div>
      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Имя</TableHead>
                <TableHead>Telegram</TableHead>
                <TableHead>Должность</TableHead>
                <TableHead>Опыт</TableHead>
                <TableHead>Теги</TableHead>
                <TableHead class="text-right" />
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="mentorService.items.value.total === 0" class="h-24">
                <TableCell colspan="7" class="text-center">
                  Менторы не найдены
                </TableCell>
              </TableRow>
              <TableRow v-for="mentor in mentorService.items.value.items" :key="mentor.id">
                <TableCell>{{ mentor.name }}</TableCell>
                <TableCell>{{ mentor.tg }}</TableCell>
                <TableCell>{{ mentor.occupation }}</TableCell>
                <TableCell>{{ mentor.experience }}</TableCell>
                <TableCell>
                  <div class="flex flex-wrap gap-1">
                    <span
                      v-for="tag in mentor.profTags"
                      :key="tag.id"
                      class="px-2 py-1 text-xs rounded-full bg-primary/10"
                    >
                      {{ tag.title }}
                    </span>
                  </div>
                </TableCell>
                <TableCell class="text-right">
                  <Button variant="ghost" size="sm" @click="handleEditMentor(mentor)">
                    <Pencil class="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" class="text-destructive" @click="mentorService.delete(mentor.id)">
                    <Trash class="h-4 w-4" />
                  </Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <div class="mt-4 flex justify-end">
        <Pagination v-slot="{ page }" :items-per-page="10" :total="mentorService.items.value.total" :sibling-count="1" show-edges :default-page="1">
          <PaginationList v-slot="{ items }" class="flex items-center gap-1">
            <PaginationFirst />
            <PaginationPrev />

            <template v-for="(item, index) in items">
              <PaginationListItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                <Button class="w-10 h-10 p-0" :variant="item.value === page ? 'default' : 'outline'" @click="mentorService.changePagination(item.value)">
                  {{ item.value }}
                </Button>
              </PaginationListItem>
              <PaginationEllipsis v-else :key="item.type" :index="index" />
            </template>

            <PaginationNext />
            <PaginationLast />
          </PaginationList>
        </Pagination>
      </div>

      <MentorModal
        v-model:is-open="isModalOpen"
        :mentor-id="selectedMentor?.id!"
        :member-id="selectedMentor?.memberId!"
        @saved="mentorService.search"
      />
    </div>
  </AdminLayout>
</template>
