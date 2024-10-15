<template>
  <div v-if="setting">
    <div class="border-solid border-b border-gray-300 pb-2">
      <div class="text-xl font-bold mb-4">基本信息</div>
      <div class="text-sm text-gray-500 mb-4">在这里设置网站的基本信息</div>
      <el-form :model="setting" label-width="90px" class="flex flex-col md:flex-row">
        <div class="flex-1 mb-6 md:mb-0 md:mr-6">
          <el-form-item label="博客问候">
            <template #label>
              <span class="">博客问候</span>
              <el-tooltip content="首页第一行大标题" placement="top">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </template>
            <el-input v-model="setting.greeting" />
          </el-form-item>
          <el-form-item label="角色">
            <template #label>
              <span>角色</span>
              <el-tooltip content="首页第二行标题" placement="top">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </template>
            <el-input v-model="setting.role" />
          </el-form-item>
          <el-form-item label="简介">
            <template #label>
              <span>简介</span>
              <el-tooltip content="首页第三行介绍部分" placement="top">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </template>
            <el-input v-model="setting.philosophy" />
          </el-form-item>
        </div>
        <div class="flex-1">
          <el-form-item label="博客名称">
            <template #label>
              <span>博客名称</span>
              <el-tooltip content="博客名称，仅适用于前台" placement="top">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </template>
            <el-input v-model="setting.name" />
          </el-form-item>
          <el-form-item label="博客地址">
            <template #label>
              <span>博客地址</span>
              <el-tooltip content="博客的前台部署地址" placement="top">
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </template>
            <el-input v-model="setting.address" />
          </el-form-item>
          <el-form-item label="博客头像">
            <template #label>
              <span>博客头像</span>
              <el-tooltip
                content="博客头像，仅适用于前台，前台看到的几乎所有头像都是这个"
                placement="top"
              >
                <el-icon><InfoFilled /></el-icon>
              </el-tooltip>
            </template>
            <el-input v-model="setting.avatar" />
          </el-form-item>
        </div>
        <div class="flex flex-row justify-end gap-4">
          <el-button type="primary" @click="save(1)">保存</el-button>
        </div>
      </el-form>
    </div>
    <div class="flex flex-col mt-4 w-full">
      <div class="border-solid border-b border-gray-300 pb-2">
        <div class="text-xl font-bold mb-4">链接徽章</div>
        <div class="text-sm text-gray-500 mb-4">首页的外部链接，点击会跳转至对应链接</div>
        <div class="text-sm text-gray-500 mb-4">新增后需要点击保存以生效</div>
        <div v-for="item in setting.badge" :key="item.index">
          <el-form class="flex flex-col md:flex-row md:gap-4">
            <div class="flex flex-row gap-4 flex-1">
              <el-form-item label="标题">
                <el-input v-model="item.title" />
              </el-form-item>
              <el-form-item label="链接" class="flex-1">
                <el-input v-model="item.url" />
              </el-form-item>
            </div>
            <div class="flex flex-row gap-4 flex-1">
              <el-form-item label="图片" class="flex-1">
                <el-input v-model="item.img" />
              </el-form-item>
              <el-form-item label="排序" class="w-28">
                <el-input-number v-model="item.orderNum" :min="0" controls-position="right" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="deleteItem(2, item.index)">删除</el-button>
              </el-form-item>
            </div>
          </el-form>
        </div>
        <div class="flex flex-row justify-end gap-4">
          <el-button type="primary" @click="add(2)">新增</el-button>
          <el-button type="primary" @click="save(2)">保存</el-button>
        </div>
      </div>
      <div class="border-solid border-b border-gray-300 mt-4 pb-2">
        <div class="text-xl font-bold mb-4">底部链接</div>
        <div class="text-sm text-gray-500 mb-4">页脚信息，留空 url 则不会跳转，作为纯文本显示</div>
        <div class="text-sm text-gray-500 mb-4">
          上排仅接受站内链接（如/article/1），下排仅接受站外链接（如https://github.com）
        </div>
        <div v-for="item in setting.footer" :key="item.index">
          <el-form class="flex flex-col md:flex-row md:gap-4">
            <div class="flex flex-row gap-4 flex-1">
              <el-form-item label="标题" class="flex-1">
                <el-input v-model="item.text" />
              </el-form-item>
              <el-form-item label="位置" class="flex-1 w-28">
                <el-select v-model="item.position" placeholder="位置">
                  <el-option label="上排" :value="1" :key="1" />
                  <el-option label="下排" :value="2" :key="2" />
                </el-select>
              </el-form-item>
            </div>
            <div class="flex flex-row gap-4 flex-1">
              <el-form-item label="url" class="flex-1">
                <el-input v-model="item.url" />
              </el-form-item>
              <el-form-item label="排序" class="w-28">
                <el-input-number v-model="item.orderNum" :min="0" controls-position="right" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="deleteItem(3, item.index)">删除</el-button>
              </el-form-item>
            </div>
          </el-form>
        </div>
        <div class="flex flex-row justify-end gap-4">
          <el-button type="primary" @click="add(3)">新增</el-button>
          <el-button type="primary" @click="save(3)">保存</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reqSettingDelete, reqSettingList, reqSettingUpdate } from '@/api/setting'
import type { SettingData } from '@/api/setting/type'
import { ref, onMounted } from 'vue'
import { InfoFilled } from '@element-plus/icons-vue'
const setting = ref<SettingData>({})

const getSetting = async () => {
  const res = await reqSettingList()
  setting.value = res.data
}

const add = (type: number) => {
  switch (type) {
    case 2:
      setting.value.badge!.push({})
      break
    case 3:
      setting.value.footer!.push({})
      break
  }
}

const save = (type: number) => {
  switch (type) {
    case 1:
      reqSettingUpdate({
        greeting: setting.value.greeting,
        role: setting.value.role,
        philosophy: setting.value.philosophy,
        name: setting.value.name,
        address: setting.value.address,
        avatar: setting.value.avatar,
      })
      break
    case 2:
      setting.value.badge = setting.value.badge!.map((item) => {
        return {
          index: item.index ?? item.title,
          title: item.title,
          url: item.url,
          img: item.img,
          orderNum: item.orderNum,
        }
      })
      reqSettingUpdate({ badge: setting.value.badge })
      break
    case 3:
      setting.value.footer = setting.value.footer!.map((item) => {
        return {
          index: item.index ?? item.text,
          text: item.text,
          url: item.url,
          position: item.position,
          orderNum: item.orderNum,
        }
      })
      reqSettingUpdate({ footer: setting.value.footer })
      break
  }
}

const deleteItem = (type: number, index?: string) => {
  switch (type) {
    case 2:
      reqSettingDelete(type, index!)
      break
    case 3:
      reqSettingDelete(type, index!)
      break
  }
}

onMounted(() => {
  getSetting()
})
</script>

<style scoped></style>
