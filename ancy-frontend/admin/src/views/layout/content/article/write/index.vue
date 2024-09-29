<template>
  <div class="app-container">
    <el-form ref="form" :model="article" label-width="90px">
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="文章标题" prop="title">
            <el-input
              v-model="article.title"
              placeholder="请输入文章标题"
              maxlength="30"
            />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="分类">
            <el-select v-model="article.categoryId" placeholder="请选择">
              <el-option
                v-for="category in categoryList"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="标签">
            <el-select v-model="article.tags" placeholder="请选择" multiple>
              <el-option
                v-for="tag in tagList"
                :key="tag.id"
                :label="tag.name"
                :value="tag.id"
              />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="文章摘要">
            <el-input v-model="article.summary" type="textarea" />
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="允许评论">
            <el-radio-group v-model="article.isComment">
              <el-radio :key="'0'" :label="'0'">正常</el-radio>
              <el-radio :key="'1'" :label="'1'">停用</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
        <el-col :span="6">
          <el-form-item label="是否置顶">
            <el-radio-group v-model="article.isTop">
              <el-radio :key="'0'" :label="'0'">是</el-radio>
              <el-radio :key="'1'" :label="'1'">否</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-col>
      </el-row>
      <el-row :gutter="20" />

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="缩略图">
            <el-upload
              :file-list="article.thumbnail"
              class="upload-demo"
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
                将文件拖到此处，或
                <em>点击上传</em>
              </div>
              <template v-slot:tip>
                <div class="el-upload__tip">
                  只能上传jpg/png文件，且不超过500kb
                </div>
              </template>
            </el-upload>
          </el-form-item>
        </el-col>
        <el-col :span="12">
          <el-form-item>
            <el-button type="primary" size="medium" @click="handleSubmit">
              {{ aId ? '更新' : '发布' }}
            </el-button>
          </el-form-item>
          <el-form-item>
            <el-button v-if="!aId" type="info" @click="handleSave">
              保存到草稿箱
            </el-button>
          </el-form-item>
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

//   const res = await Promise.all(
//   files.map((file) => {
//     return new Promise((rev, rej) => {
//       const form = new FormData();
//       form.append('file', file);

//       axios
//         .post('/api/img/upload', form, {
//           headers: {
//             'Content-Type': 'multipart/form-data'
//           }
//         })
//         .then((res) => rev(res))
//         .catch((error) => rej(error));
//     });
//   })
// );

// callback(res.map((item) => item.data.url));

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
      cancelButtonText: '取消',
    }).then(() => {
      article.value = JSON.parse(localStorage.getItem('write-content')!)
    })
  }
})
</script>

<style scoped lang="scss"></style>
