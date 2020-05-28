<template>
    <div class="frame">
        <slot />
    </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";

@Component
export default class Frame extends Vue {
}
</script>

<style lang="scss">
.frame {
    background-color: $color-white;
    padding: 2rem;
}
</style>
