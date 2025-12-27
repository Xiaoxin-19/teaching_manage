import { inject, getCurrentInstance } from 'vue'
// 假设组件路径如下，请根据实际情况调整引用
import type { ConfirmOptions } from '../components/GlobalConfirmDialog.vue'

// 定义核心函数的类型签名
export type ShowConfirmFunction = (
  title: string,
  message: string,
  options?: ConfirmOptions
) => Promise<boolean>

// 全局变量作为 fallback
let globalShowConfirm: ShowConfirmFunction | null = null

// 提供给 App.vue 注册使用
export function registerConfirm(fn: ShowConfirmFunction) {
  globalShowConfirm = fn
}

export function useConfirm() {
  // 1. 尝试从组件注入
  let showConfirm: ShowConfirmFunction | undefined | null = null

  if (getCurrentInstance()) {
    showConfirm = inject<ShowConfirmFunction>('showConfirm')
  }

  // 2. 如果注入失败，尝试使用全局注册的函数
  const confirmFn = showConfirm || globalShowConfirm

  if (!confirmFn) {
    throw new Error(
      'useConfirm() 必须在提供了 showConfirm 的组件树中使用，或者通过 registerConfirm 在 App.vue 中注册。'
    )
  }

  /**
   * 基础调用
   */
  const confirm = (title: string, message: string, options?: ConfirmOptions) => {
    return confirmFn(title, message, options)
  }

  /**
   * 快捷方法：删除确认 (红色危险样式)
   * @param itemName 要删除的项目名称
   */
  const confirmDelete = (itemName: string) => {
    return confirmFn(
      '删除确认',
      `确定要永久删除 "${itemName}" 吗？此操作无法撤销。`,
      {
        type: 'error',
        confirmText: '删除',
        cancelText: '取消'
      }
    )
  }

  /**
   * 快捷方法：警告确认 (黄色警告样式)
   * @param message 警告信息
   */
  const confirmWarning = (message: string) => {
    return confirmFn('警告', message, {
      type: 'warning',
      confirmText: '继续'
    })
  }

  /**
   * 快捷方法：信息确认 (蓝色信息样式)
   * @param message 提示信息
   */
  const confirmInfo = (message: string) => {
    return confirmFn('提示', message, {
      type: 'info',
      confirmText: '确定'
    })
  }

  return {
    confirm,
    confirmDelete,
    confirmWarning,
    confirmInfo
  }
}