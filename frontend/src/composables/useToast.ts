import { inject, getCurrentInstance } from 'vue'

// 定义类型以获得良好的代码提示
export type ToastType = 'success' | 'error' | 'info' | 'warning'
export type ToastLocation = 'top-right' | 'top-center' | 'bottom-right' | 'bottom-center'

export type ShowToastFunction = (
  msg: string, 
  type?: ToastType, 
  location?: ToastLocation, 
  timeout?: number
) => void

// 全局变量作为 fallback，解决非组件环境下无法 inject 的问题
let globalShowToast: ShowToastFunction | null = null

// 提供给 App.vue 注册使用
export function registerToast(fn: ShowToastFunction) {
  globalShowToast = fn
}

export function useToast() {
  // 优先尝试 inject，如果失败则使用全局变量
  let showToast: ShowToastFunction | null = null
  
  if (getCurrentInstance()) {
    showToast = inject('showToast', null)
  }
  
  showToast = showToast || globalShowToast

  if (!showToast) {
    throw new Error('useToast必须在App.vue提供了showToast的环境下使用，或者通过 registerToast 注册')
  }

  const toast = showToast

  // 返回封装好的便捷方法
  return {
    // 原始通用方法
    open: toast,
    // 快捷方法
    success: (msg: string, loc?: ToastLocation) => toast(msg, 'success', loc),
    error: (msg: string, loc?: ToastLocation) => toast(msg, 'error', loc),
    info: (msg: string, loc?: ToastLocation) => toast(msg, 'info', loc),
    warning: (msg: string, loc?: ToastLocation) => toast(msg, 'warning', loc),
  }
}