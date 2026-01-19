import '@mdi/font/css/materialdesignicons.css' // Ensure you are using css-loader
import {createApp, ref,  provide} from 'vue'
import App from './App.vue'
import { router } from './router'


// Vuetify
import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import { zhHans } from 'vuetify/locale'


const vuetify = createVuetify({
  components,
  directives,
  locale: {
    locale: 'zhHans',
    messages: { zhHans },
  },
  icons: {
    defaultSet: 'mdi', // This is already the default value - only for display purposes
  },
  defaults: {
    VTextField: {
      autocomplete: 'off',
    },
    VTextarea: {
      autocomplete: 'off',
    },
    VAutocomplete: {
      autocomplete: 'off',
    },
    VCombobox: {
      autocomplete: 'off',
    },
    VSelect: {
      autocomplete: 'off',
    },
  },
})


createApp(App).use(vuetify).use(router).mount('#app')


