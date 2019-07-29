<template>
  <div class="main-view">
    <div>{{message}}</div>
    <hr>
    <SampleComponent></SampleComponent>
    <hr>
    <SampleComponent></SampleComponent>

  </div>
</template>

<script>
import Astor from '@/mixins/Astor.js'
import Activator  from '@/mixins/Activator.js'

import SampleComponent from '@/components/Sample'

export default {
  name: 'main',
  mixins: [Astor, Activator],
  components: {
    SampleComponent
  },
  data() {
    return {
      message: "This is the main page"
    }
  },
  methods: {
    onViewActivated: function(payload) {
      console.log("onVieActivated")
    },
    onViewDeactivated: function(payload) {
      console.log("onVieActivated")
    },
    onSampleMessage1: function(payload) {
      console.log("onSampleMessage1")
      this.message = payload
    },
  },
  mounted: function() {
    this.__astor().listen('sample.message1', this.onSampleMessage1);
  }
}
</script>

<style scoped lang="scss">
  @import "~@/styles/views/_MainView.scss";
</style>
