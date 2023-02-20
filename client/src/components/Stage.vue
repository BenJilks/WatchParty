<script setup lang="ts">
</script>

<template>
  <div class="floor"></div>
  <div class="floor-overlay floor-back"></div>
  <div class="floor-overlay floor-sides"></div>
  <img src="/cinema/curtain_left.png" class="curtain curtain_left" draggable="false" />
  <img src="/cinema/curtain_right.png" class="curtain curtain_right" draggable="false" />
</template>

<style scoped>
  @media (max-aspect-ratio: 4/5) {
    .curtain {
      --curtain-offset: 100vw;
    }
  }

  .curtain {
    position: absolute;
    width: auto;
    height: calc(100% - var(--seat-height) - 1vh);
    top: 0;
  }

  .curtain_left {
    left: var(--curtain-offset);
  }

  .curtain_right {
    right: var(--curtain-offset);
  }

  .floor {
    position: absolute;
    height: var(--floor-height);
    width: 200%;

    background-repeat: repeat-x;
    background-size: var(--floor-height);
    background-position-x: 50%;
    background-image: url("/cinema/floor.jpg");
    transform:
        perspective(50vh)
        rotateX(80deg);

    left: -50%;
    bottom: calc(var(--seat-height) - var(--floor-height)*0.39);
  }

  .floor-overlay {
    position: absolute;
    user-select: none;
    pointer-events: none;

    width: 100%;
    height: calc(var(--floor-height)*0.25);
    bottom: var(--seat-height);
  }

  .floor-sides {
    background: linear-gradient(90deg,
    rgba(0,0,0,0.9) 0%,
    rgba(0,0,0,0.2) 50%,
    rgba(0,0,0,0.9) 100%);
  }

  .floor-back {
    background: linear-gradient(0deg,
    rgba(0,0,0,0) 0%,
    rgba(0,0,0,1) 100%);
  }
</style>
