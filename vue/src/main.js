import Vue from 'vue'
import App from '@/App.vue'
import Vuetify from 'vuetify'
import "vue-material-design-icons/styles.css"
import VTooltip from 'v-tooltip'
import Astor  from '@/Astor.js'

Vue.config.productionTip = false

Vue.use(VTooltip);
Vue.use(Vuetify);

Astor.init();

document.addEventListener('astilectron-ready', function() {
  new Vue({
    render: h => h(App),
  }).$mount('#vue');
});


