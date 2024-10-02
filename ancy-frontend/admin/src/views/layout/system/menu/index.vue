<template>
  <div class="app-container">
    <el-form ref="queryForm" :model="queryParams" size="small" class="query-form">
      <el-input
        v-model="queryParams.name"
        placeholder="请输入菜单名称"
        clearable
        @keyup.enter="getMenuPage()"
      />
      <el-select v-model="queryParams.status" placeholder="菜单状态" clearable size="small">
        <el-option key="0" label="正常" value="0" />
        <el-option key="1" label="停用" value="1" />
      </el-select>
      <el-button type="primary" :icon="Search" @click="getMenuPage()">搜索</el-button>
    </el-form>

    <el-row :gutter="10">
      <el-col :span="1.5">
        <el-button type="primary" plain :icon="Plus" @click="handleAdd">新增</el-button>
      </el-col>
    </el-row>

    <el-table
      :data="menuTree"
      :default-expand-all="false"
      v-loading="loading"
      row-key="id"
      :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
    >
      <el-table-column
        prop="menuName"
        label="菜单名"
        :show-overflow-tooltip="true"
        width="160"
        align="center"
      />
      <el-table-column prop="icon" label="图标" align="center" width="60">
        <template v-slot="scope">
          <el-icon>
            <component :is="iconComponent(scope.row.icon)" />
          </el-icon>
        </template>
      </el-table-column>
      <el-table-column prop="orderNum" label="排序" align="center" width="60" />
      <el-table-column prop="menuType" label="类型" align="center" width="60">
        <template v-slot="scope">
          <span v-if="scope.row.menuType === 'M'">目录</span>
          <span v-if="scope.row.menuType === 'C'">菜单</span>
          <span v-if="scope.row.menuType === 'F'">按钮</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="perms"
        label="权限标识"
        :show-overflow-tooltip="true"
        width="160"
        align="center"
      />
      <el-table-column
        prop="component"
        label="组件路径"
        :show-overflow-tooltip="true"
        width="160"
        align="center"
      />
      <el-table-column prop="status" label="启用" width="80" align="center">
        <template v-slot="scope">
          <el-switch
            v-model="scope.row.status"
            active-value="0"
            inactive-value="1"
            :loading="statusLoading"
            @change="handleStatusChange(scope.row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="创建时间" align="center" prop="createTime" width="100">
        <template v-slot="scope">
          {{ scope.row.createTime }}
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" fixed="right" width="240">
        <template v-slot="scope">
          <el-button size="small" type="text" :icon="Edit" @click="handleUpdate(scope.row)">
            修改
          </el-button>
          <el-button
            size="small"
            type="text"
            :icon="Plus"
            @click="handleAdd(scope.row)"
            :disabled="scope.row.menuType === 'F'"
          >
            新增
          </el-button>
          <el-button size="small" type="text" :icon="Delete" @click="handleDelete(scope.row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加或修改菜单对话框 -->
    <el-dialog :title="title" v-model="open" class="dialog-form">
      <el-form ref="form" :model="menu" :rules="rules" label-width="auto" label-position="top">
        <el-row>
          <el-col :span="24">
            <el-form-item label="上级菜单">
              <el-tree-select
                v-model="menu.parentId"
                :data="menuOptions"
                placeholder="选择上级菜单"
                check-strictly
                :props="{ children: 'children', label: 'menuName', value: 'id' }"
              />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="菜单类型" prop="menuType">
              <el-radio-group v-model="menu.menuType">
                <el-radio label="M">目录</el-radio>
                <el-radio label="C">菜单</el-radio>
                <el-radio label="F">按钮</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item v-if="menu.menuType != 'F'" label="菜单图标">
              <el-select-v2
                filterable
                v-model="menu.icon"
                :options="icons"
                placeholder="选择图标"
                style="width: 240px"
              >
                <template #default="{ item }">
                  <el-icon>
                    <component :is="iconComponent(item.value)" />
                  </el-icon>
                  <span>{{ item.label }}</span>
                </template>
              </el-select-v2>
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item label="菜单名称" prop="menuName">
              <el-input v-model="menu.menuName" placeholder="菜单名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="显示排序" prop="orderNum">
              <el-input-number v-model="menu.orderNum" controls-position="right" :min="0" />
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item v-if="menu.menuType != 'F'" prop="path">
              <template v-slot:label>
                <span>
                  <el-tooltip content="访问的路由地址，如：`user`" placement="top">
                    路由地址
                  </el-tooltip>
                </span>
              </template>
              <el-input v-model="menu.path" placeholder="路由地址" />
            </el-form-item>
          </el-col>

          <el-col v-if="menu.menuType == 'C'" :span="12">
            <el-form-item prop="component">
              <template v-slot:label>
                <span>
                  <el-tooltip
                    content="访问的组件路径，如：`system/user/index`，默认在`views`目录下"
                    placement="top"
                  >
                    组件路径
                  </el-tooltip>
                </span>
              </template>
              <el-input v-model="menu.component" placeholder="组件路径" />
            </el-form-item>
          </el-col>

          <el-col :span="12" v-if="menu.menuType != 'M'">
            <el-form-item>
              <el-input v-model="menu.perms" placeholder="权限标识" maxlength="100" />
              <template v-slot:label>
                <span>
                  <el-tooltip
                    content="控制器中定义的权限字符，如：@PreAuthorize(`@ss.hasPermi('system:user:list')`)"
                    placement="top"
                  >
                    权限字符
                  </el-tooltip>
                </span>
              </template>
            </el-form-item>
          </el-col>

          <el-col :span="12">
            <el-form-item v-if="menu.menuType != 'F'">
              <template v-slot:label>
                <span>
                  <el-tooltip
                    content="选择隐藏则路由将不会出现在侧边栏，但仍然可以访问"
                    placement="top"
                  >
                    显示状态
                  </el-tooltip>
                </span>
              </template>
              <el-radio-group v-model="menu.visible">
                <el-radio :key="'0'" :label="'0'">显示</el-radio>
                <el-radio :key="'1'" :label="'1'">隐藏</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item v-if="menu.menuType != 'F'">
              <template v-slot:label>
                <span>
                  <el-tooltip
                    content="选择停用则路由将不会出现在侧边栏，也不能被访问"
                    placement="top"
                  >
                    菜单状态
                  </el-tooltip>
                </span>
              </template>
              <el-radio-group v-model="menu.status">
                <el-radio key="0" label="0">正常</el-radio>
                <el-radio key="1" label="1">停用</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <template v-slot:footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="handleSubmit">确 定</el-button>
          <el-button @click="handleCancel">取 消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {
  reqMenuAdd,
  reqMenuDelete,
  reqMenuList,
  reqMenuTree,
  reqMenuUpdate,
} from '@/api/system/menu'
import type { MenuAddParams, MenuListData, MenuPageParams } from '@/api/system/menu/type'
import { onMounted, ref, watch } from 'vue'
import { Search, Plus, Delete, Edit } from '@element-plus/icons-vue'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import type { FormRules } from 'element-plus'
import { toggleStatus } from '@/utils/toggleStatus'

const queryParams = ref<MenuPageParams>({})
const menuList = ref<MenuListData[]>([])
const loading = ref<boolean>(false)
const statusLoading = ref<boolean>(false)
const open = ref<boolean>(false)
const title = ref<string>('')
const menu = ref<MenuAddParams>({})
const rules = ref<FormRules>({})
const menuOptions = ref<MenuListData[]>([])
const icons = ref<{ value: string; label: string }[]>([])
const menuTree = ref<MenuListData[]>([])

watch(
  () => menuList.value,
  () => {
    menuTree.value = buildTree(menuList.value)
  },
)

const buildTree = (data: MenuListData[]) => {
  const tree: MenuListData[] = []
  const map = new Map()
  data.forEach((item) => {
    map.set(item.id, item)
  })
  data.forEach((item) => {
    if (item.parentId) {
      const parent = map.get(item.parentId)
      parent.children = parent.children || []
      parent.children.push(item)
    } else {
      tree.push(item)
    }
  })
  return tree
}

const iconComponent = (icon: string) => {
  return ElementPlusIconsVue[icon as keyof typeof ElementPlusIconsVue]
}

const getMenuPage = async (page: number = 1) => {
  loading.value = true
  queryParams.value.pageNum = page
  const res = await reqMenuList(queryParams.value)
  menuList.value = res.data
  loading.value = false
}

const handleStatusChange = async (row: MenuListData) => {
  statusLoading.value = true
  try {
    await reqMenuUpdate(row.id, { status: row.status })
  } catch {
    toggleStatus(row)
  } finally {
    statusLoading.value = false
  }
}

const getMenuOptions = async () => {
  const res = await reqMenuTree()
  menuOptions.value = res.data
}

const getIcons = () => {
  icons.value = Object.keys(ElementPlusIconsVue).map((item) => {
    return {
      value: item,
      label: item,
    }
  })
}

const handleAdd = (row?: MenuListData) => {
  menu.value = {}
  if (row) {
    menu.value.parentId = row.id
  }
  open.value = true
  title.value = '新增'
}

const handleUpdate = (row: MenuListData) => {
  menu.value = row
  open.value = true
  title.value = '修改'
}

const handleDelete = (row: MenuListData) => {
  ElMessageBox.confirm('删除?', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(async () => {
    await reqMenuDelete(row.id)
  })
}

const handleSubmit = async () => {
  if (menu.value.id) {
    await reqMenuUpdate(menu.value.id, menu.value)
  } else {
    await reqMenuAdd(menu.value)
    ElMessage.success('新增成功')
  }
  open.value = false
  getMenuPage()
}

const handleCancel = () => {
  menu.value = {}
  open.value = false
}

onMounted(() => {
  getMenuPage()
  getMenuOptions()
  getIcons()
})
</script>

<style scoped lang="scss"></style>
