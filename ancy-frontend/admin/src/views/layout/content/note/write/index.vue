<template>
  <div class="app-container">
    <el-row>
      <MdEditor v-model="note.content" />
    </el-row>
    <el-form ref="form" :model="note" label-width="90px">
      <el-row :gutter="20" style="margin: 10px 0">
        <el-col :span="7">
          <el-form-item label="评论" label-width="40px">
            <el-radio-group v-model="note.isComment">
              <el-radio key="1" label="1">允许</el-radio>
              <el-radio key="0" label="0">禁止</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="7" style="text-align: right">
          <el-form-item label="置顶" label-width="40px">
            <el-radio-group v-model="note.isTop">
              <el-radio key="1" label="1">是</el-radio>
              <el-radio key="0" label="0">否</el-radio>
              {{ note.isTop }}
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="排序" label-width="40px">
            <el-input-number v-model="note.orderNum" controls-position="right" :min="0" />
          </el-form-item>
        </el-col>
        <el-col :span="10" style="text-align: right">
          <el-button type="primary" @click="handleSubmit">
            {{ aId ? '更新' : '发布' }}
          </el-button>
        </el-col>
      </el-row>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import MdEditor from '@/components/MdEditor.vue'
import { useRoute, type LocationQueryValue } from 'vue-router'
import type { NoteAddParams } from '@/api/note/type'
import { reqNoteAdd, reqNoteGetById, reqNoteUpdate } from '@/api/note'

const route = useRoute()
const note = ref<NoteAddParams>({})
const aId = ref<number | null>(null)

watch(
  () => route.query.id,
  async (newVal: LocationQueryValue | LocationQueryValue[]) => {
    aId.value = Number(newVal)
    if (aId.value) {
      await getNote(aId.value)
    } else {
      note.value = {}
    }
  },
  {
    immediate: true,
  },
)

watch(
  note,
  (newVal) => {
    localStorage.setItem('write-content', JSON.stringify(newVal))
  },
  {
    deep: true,
  },
)

const getNote = async (id: number) => {
  const res = await reqNoteGetById(id)
  note.value = {
    ...res.data,
  }
}

const handleSubmit = async () => {
  note.value.status = '0'
  if (aId.value) {
    await reqNoteUpdate(aId.value, note.value)
  } else {
    await reqNoteAdd(note.value)
    ElMessage.success('发布成功')
  }
}

onMounted(async () => {
  if (aId.value) {
    await getNote(aId.value)
  }

  if (!aId.value && localStorage.getItem('write-content')) {
    ElMessageBox.confirm('是否恢复上次编辑内容', '提示', {
      confirmButtonText: '恢复',
      cancelButtonText: '删除本地内容',
    })
      .then(() => {
        note.value = JSON.parse(localStorage.getItem('write-content')!)
      })
      .catch(() => {
        localStorage.removeItem('write-content')
      })
  }
})
</script>

<style scoped lang="scss">
.img-box {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  border: 1px solid $ac-primary-color;
  border-radius: 5px;

  img {
    width: 100%;
  }

  .el-button {
    position: absolute;
    top: 0;
    right: 0;
    display: none;
  }

  &:hover {
    .el-button {
      display: block;
    }
  }
}
</style>
