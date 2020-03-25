<template>
  <div class="notification__container">
    <Bell class="notification__icon" />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Bell from "../../assets/images/icons/bell-solid.svg";

@Component({
  components: { Bell }
})
export default class Notification extends Vue {}
</script>

<style lang="scss">
.notification {
  &__container {
    padding: 2.15rem;
    cursor: pointer;
  }
  &__icon {
    display: inline-block;
    position: relative;
    top: 0.25rem;
    height: 2rem;
    color: $color-grey-1;
    transition: all 0.2s;
  }

  &__container:hover {
    > .notification__icon {
      color: $color-complementary;
    }
  }

  &__container:active {
    > .notification__icon {
      color: $color-complementary-dark;
    }
  }
}
</style>
