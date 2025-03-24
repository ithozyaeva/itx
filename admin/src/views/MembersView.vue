<script setup lang="ts">
import AdminLayout from '@/components/layout/AdminLayout.vue'
import MemberModal from '@/components/modals/MemberModal.vue'
import MentorModal from '@/components/modals/MentorModal.vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Pagination, PaginationEllipsis, PaginationFirst, PaginationLast, PaginationList, PaginationListItem, PaginationNext, PaginationPrev } from '@/components/ui/pagination'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { memberService } from '@/services/memberService'
import { onMounted, ref } from 'vue'

import Pencil from '~icons/lucide/pencil'
import Plus from '~icons/lucide/plus'
import Trash from '~icons/lucide/trash'

const isModalOpen = ref(false)
const currentMemberId = ref<number | null>(null)

const isMentorModalOpen = ref(false)
const selectedMemberId = ref<number | null>(null)

/**
 * Обработчик открытия модального окна для одобавления участника.
 */
function openAddModal() {
  currentMemberId.value = null
  isModalOpen.value = true
}

/**
 * Обработчик открытия модального окна для редактирования участника.
 *
 * @param id - ID участника.
 */
function openEditModal(id: number) {
  currentMemberId.value = id
  isModalOpen.value = true
}

/**
 * Обработчик открытия модального окна для создания ментора из участинка.
 *
 * @param memberId - ID участника.
 */
function handleMakeMentor(memberId: number) {
  selectedMemberId.value = memberId
  isMentorModalOpen.value = true
}

onMounted(memberService.search)
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <div class="flex justify-between items-center">
        <h1 class="text-3xl font-bold">
          Участники сообщества
        </h1>
        <Button @click="openAddModal">
          <Plus class="mr-2 h-4 w-4" />
          Добавить участника
        </Button>
      </div>

      <Card>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Имя</TableHead>
                <TableHead>Telegram</TableHead>
                <TableHead />
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-if="memberService.items.value.total === 0" class="h-24">
                <TableCell colspan="7" class="text-center">
                  Участники не найдены
                </TableCell>
              </TableRow>
              <TableRow v-for="member in memberService.items.value.items" :key="member.id">
                <TableCell>{{ member.name }}</TableCell>
                <TableCell>{{ member.tg }}</TableCell>
                <TableCell class="text-right">
                  <Button variant="ghost" size="sm" @click="openEditModal(member.id)">
                    <Pencil class="h-4 w-4" />
                  </Button>
                  <Button variant="ghost" size="sm" class="text-destructive" @click="memberService.delete(member.id)">
                    <Trash class="h-4 w-4" />
                  </Button>
                  <Button v-if="!member.isMentor" variant="ghost" size="sm" class="text-primary" @click="handleMakeMentor(member.id)">
                    Сделать ментором
                  </Button>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>
      <div class="mt-4 flex justify-end">
        <Pagination v-slot="{ page }" :items-per-page="10" :total="memberService.items.value.total" :sibling-count="1" show-edges :default-page="1">
          <PaginationList v-slot="{ items }" class="flex items-center gap-1">
            <PaginationFirst />
            <PaginationPrev />

            <template v-for="(item, index) in items">
              <PaginationListItem v-if="item.type === 'page'" :key="index" :value="item.value" as-child>
                <Button class="w-10 h-10 p-0" :variant="item.value === page ? 'default' : 'outline'" @click="memberService.changePagination(item.value)">
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
    </div>

    <MemberModal
      v-model:is-open="isModalOpen"
      :member-id="currentMemberId"
      @saved="memberService.search"
    />

    <!-- Модальное окно для создания ментора -->
    <MentorModal
      v-model:is-open="isMentorModalOpen"
      :member-id="selectedMemberId!"
      @saved="memberService.search"
    />
  </AdminLayout>
</template>
