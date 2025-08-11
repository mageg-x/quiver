import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index'

import PrimeVue from 'primevue/config';
import Aura from '@primeuix/themes/aura';
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'

// PrimeVue Components
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Card from 'primevue/card'
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dialog from 'primevue/dialog'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Menu from 'primevue/menu'
import Menubar from 'primevue/menubar'
import Sidebar from 'primevue/sidebar'
import Tree from 'primevue/tree'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import Toolbar from 'primevue/toolbar'
import SplitButton from 'primevue/splitbutton'
import Badge from 'primevue/badge'
import Tag from 'primevue/tag'
import ProgressSpinner from 'primevue/progressspinner'
import Skeleton from 'primevue/skeleton'

// Styles
import './style.css'
import 'primeicons/primeicons.css'

const app = createApp(App)
app.use(router)
app.use(PrimeVue, {
    theme: {
        preset: Aura,
        ripple: true,
        options: {
            prefix: 'p',
            darkModeSelector: 'system',
            cssLayer: false
        }
    }
});
app.use(ToastService)
app.use(ConfirmationService)

// Register PrimeVue components
app.component('Button', Button)
app.component('InputText', InputText)
app.component('Password', Password)
app.component('Card', Card)
app.component('Toast', Toast)
app.component('ConfirmDialog', ConfirmDialog)
app.component('DataTable', DataTable)
app.component('Column', Column)
app.component('Dialog', Dialog)
app.component('Textarea', Textarea)
app.component('Dropdown', Dropdown)
app.component('Menu', Menu)
app.component('Menubar', Menubar)
app.component('Sidebar', Sidebar)
app.component('Tree', Tree)
app.component('TabView', TabView)
app.component('TabPanel', TabPanel)
app.component('Toolbar', Toolbar)
app.component('SplitButton', SplitButton)
app.component('Badge', Badge)
app.component('Tag', Tag)
app.component('ProgressSpinner', ProgressSpinner)
app.component('Skeleton', Skeleton)

app.mount('#app')