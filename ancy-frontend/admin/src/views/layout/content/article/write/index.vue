<template>
  <div class="app-container">
    <el-form ref="form" :model="article" label-width="90px">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-input
            v-model="article.title"
            placeholder="文章标题"
            maxlength="30"
          />
        </el-col>
        <el-col :span="6">
          <el-select v-model="article.categoryId" placeholder="选择分类">
            <el-option
              v-for="category in categoryList"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select v-model="article.tags" placeholder="选择标签" multiple>
            <el-option
              v-for="tag in tagList"
              :key="tag.id"
              :label="tag.name"
              :value="tag.id"
            />
          </el-select>
        </el-col>
      </el-row>
      <el-row :gutter="20" style="margin-top: 10px">
        <el-col :span="12">
          <el-input
            v-model="article.summary"
            type="textarea"
            placeholder="文章摘要"
          />
        </el-col>
        <el-col :span="12">
          <el-upload
            :file-list="article.thumbnail"
            list-type="picture"
            drag
            name="img"
            action="upload"
            :on-remove="fileRemove"
            :limit="1"
            :http-request="handleUpload"
            :on-exceed="onExceed"
          >
            <i class="el-icon-upload" />
            <div class="el-upload__text">
              <span>缩略图</span>
              <span>
                将文件拖到此处，或
                <em>点击上传</em>
              </span>
              <div class="el-upload__tip">上传jpg/png文件 不超过500kb</div>
            </div>
          </el-upload>
        </el-col>
      </el-row>

      <el-row :gutter="20" style="margin: 10px 0">
        <el-col :span="7">
          <el-form-item label="评论" label-width="40px">
            <el-radio-group v-model="article.isComment">
              <el-radio :key="'0'" :label="'0'">允许</el-radio>
              <el-radio :key="'1'" :label="'1'">禁止</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="7" style="text-align: right">
          <el-form-item label="置顶" label-width="40px">
            <el-radio-group v-model="article.isTop">
              <el-radio :key="'0'" :label="'0'">是</el-radio>
              <el-radio :key="'1'" :label="'1'">否</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="10" style="text-align: right">
          <el-button type="primary" size="medium" @click="handleSubmit">
            {{ aId ? '更新' : '发布' }}
          </el-button>
          <el-button v-if="aId" type="info" @click="handleSave">
            保存到草稿箱
          </el-button>
        </el-col>
      </el-row>
      <el-row>
        <MdEditor v-model="article.content" />
      </el-row>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import MdEditor from './MdEditor.vue'

const article = ref<{
  title: string
  thumbnail: string
  isTop: string
  isComment: string
  content: string
  categoryId: string
  summary: string
  tags: string[]
}>({
  title: '',
  thumbnail: '',
  isTop: '1',
  isComment: '0',
  content: '# fesdrtwewdrgfdsfgsfdgsfgd',
  categoryId: '',
  summary: '',
  tags: [],
})

watch(
  article,
  (newVal) => {
    console.log('newVal', newVal)
    localStorage.setItem('write-content', JSON.stringify(newVal))
  },
  {
    deep: true,
  },
)
const categoryList = ref([
  {
    id: 1,
    name: '前端',
  },
  {
    id: 2,
    name: '后端',
  },
  {
    id: 3,
    name: '数据库',
  },
])
const tagList = ref([
  {
    id: 1,
    name: 'Vue',
  },
  {
    id: 2,
    name: 'React',
  },
])
const aId = ref(-1)
const fileList = ref([])

const handleSubmit = (a) => {
  console.log('handleSubmit', a)
}

const handleSave = () => {
  console.log('handleSave')
}

const fileRemove = (a) => {
  console.log('fileRemove', a)
}

const onExceed = (a) => {
  console.log('onExceed', a)
}

const handleUpload = (a) => {
  console.log('handleUpload', a)
}

onMounted(() => {
  if (localStorage.getItem('write-content')) {
    ElMessageBox.confirm('是否恢复上次编辑内容', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '删除本地内容',
    })
      .then(() => {
        article.value = JSON.parse(localStorage.getItem('write-content')!)
      })
      .catch(() => {
        localStorage.removeItem('write-content')
      })
  }
})
</script>

<style scoped lang="scss"></style>
