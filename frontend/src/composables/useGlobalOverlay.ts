import { inject, getCurrentInstance } from 'vue'

export interface GlobalOverlayController {
  show: (message?: string) => void
  hide: () => void
  setMessage?: (message: string) => void
}

let globalOverlay: GlobalOverlayController | null = null

export function registerGlobalOverlay(controller: GlobalOverlayController) {
  globalOverlay = controller
}

export function useGlobalOverlay() {
  let controller: GlobalOverlayController | null = null

  if (getCurrentInstance()) {
    controller = inject('globalOverlay', null)
  }

  controller = controller || globalOverlay

  if (!controller) {
    throw new Error('useGlobalOverlay 必须在 App.vue 提供全局遮罩后使用，或通过 registerGlobalOverlay 注册')
  }

  return controller
}
