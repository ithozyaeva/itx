<script setup lang="ts">
import AdminLayout from '@/components/layout/AdminLayout.vue'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { memberService } from '@/services/memberService'
import { mentorService } from '@/services/mentorService'
import { onMounted, ref } from 'vue'

const mentorsCount = ref(0)
const membersCount = ref(0)

onMounted(async () => {
  mentorsCount.value = (await mentorService.getAll()).total
  membersCount.value = (await memberService.getAll()).total
})
</script>

<template>
  <AdminLayout>
    <div class="space-y-6">
      <h1 class="text-3xl font-bold">
        Дашборд
      </h1>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <Card>
          <CardHeader>
            <CardTitle>Менторы</CardTitle>
            <CardDescription>Общее количество менторов в системе</CardDescription>
          </CardHeader>
          <CardContent>
            <p class="text-4xl font-bold">
              {{ mentorsCount }}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Участники</CardTitle>
            <CardDescription>Общее количество участников сообщества</CardDescription>
          </CardHeader>
          <CardContent>
            <p class="text-4xl font-bold">
              {{ membersCount }}
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  </AdminLayout>
</template>
