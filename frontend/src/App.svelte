<script>
  import { onMount } from "svelte";
  import { AddEvent, DeleteEvent, GetEvents } from "../wailsjs/go/main/App";

  let currentDate = new Date();
  let selectedDate = null;

  let days = [];
  let events = [];
  let loadGeneration = 0;

  function wailsReady() {
    return typeof window !== "undefined" && window.go?.main?.App?.GetEvents;
  }

  let contextMenu = {
    visible: false,
    x: 0,
    y: 0,
    kind: "day",
    day: null,
    event: null
  };

  function formatDate(date) {
    const y = date.getFullYear();
    const m = String(date.getMonth() + 1).padStart(2, "0");
    const d = String(date.getDate()).padStart(2, "0");
    return `${y}-${m}-${d}`;
  }

  function generateCalendar(date) {
    const year = date.getFullYear();
    const month = date.getMonth();

    const firstDay = new Date(year, month, 1).getDay();
    const totalDays = new Date(year, month + 1, 0).getDate();

    let temp = [];

    for (let i = 0; i < firstDay; i++) {
      temp.push(null);
    }

    for (let d = 1; d <= totalDays; d++) {
      temp.push(new Date(year, month, d));
    }

    days = temp;
  }

  function nextMonth() {
    currentDate.setMonth(currentDate.getMonth() + 1);
    currentDate = new Date(currentDate);
    generateCalendar(currentDate);
  }

  function prevMonth() {
    currentDate.setMonth(currentDate.getMonth() - 1);
    currentDate = new Date(currentDate);
    generateCalendar(currentDate);
  }

  function selectDay(day) {
    selectedDate = day;
  }

  function monthName(date) {
    return date.toLocaleString("en-US", { month: "long", year: "numeric" });
  }

  function openDayContextMenu(e, day) {
    e.preventDefault();

    contextMenu = {
      visible: true,
      x: e.clientX,
      y: e.clientY,
      kind: "day",
      day,
      event: null
    };
  }

  function openEventContextMenu(e, eventRow) {
    e.preventDefault();
    e.stopPropagation();

    contextMenu = {
      visible: true,
      x: e.clientX,
      y: e.clientY,
      kind: "event",
      day: null,
      event: eventRow
    };
  }

  function closeContextMenu() {
    contextMenu.visible = false;
  }

  async function loadEvents() {
    const gen = loadGeneration;
    if (wailsReady()) {
      try {
        const list = await GetEvents();
        if (gen !== loadGeneration) return;
        events = list;
      } catch (e) {
        console.error(e);
        if (gen !== loadGeneration) return;
        events = [];
      }
    } else {
      events = [
        { id: "demo-1", date: "2026-04-05", title: "Meeting", color: "#4ade80" },
        { id: "demo-2", date: "2026-04-08", title: "Hack Time", color: "#60a5fa" }
      ];
    }
  }

  async function addEvent() {
    if (contextMenu.kind !== "day" || !contextMenu.day) return;
    const day = contextMenu.day;
    const title = prompt("Event adı:");
    if (!title) return;

    closeContextMenu();

    if (wailsReady()) {
      try {
        const ev = await AddEvent(formatDate(day), title, "#facc15");
        loadGeneration++;
        events = [...events, ev];
      } catch (e) {
        console.error(e);
        alert(String(e));
      }
    } else {
      const id =
        typeof crypto !== "undefined" && crypto.randomUUID
          ? crypto.randomUUID()
          : String(Date.now());
      events = [...events, { id, date: formatDate(day), title, color: "#facc15" }];
    }
  }

  async function deleteEventFromMenu() {
    if (contextMenu.kind !== "event" || !contextMenu.event) return;
    const target = contextMenu.event;
    const id = target.id;

    closeContextMenu();

    if (wailsReady()) {
      try {
        await DeleteEvent(id);
        loadGeneration++;
        events = events.filter((e) => e.id !== id);
      } catch (e) {
        console.error(e);
        alert(String(e));
      }
    } else {
      events = events.filter((e) => e.id !== id);
    }
  }

  onMount(async () => {
    generateCalendar(currentDate);
    await loadEvents();
  });
</script>

<svelte:window on:click={closeContextMenu} />

<div class="container">
  <div class="header">
    <button on:click={prevMonth}>←</button>
    <h2>{monthName(currentDate)}</h2>
    <button on:click={nextMonth}>→</button>
  </div>

  <div class="weekdays">
    <div>Sun</div>
    <div>Mon</div>
    <div>Tue</div>
    <div>Wed</div>
    <div>Thu</div>
    <div>Fri</div>
    <div>Sat</div>
  </div>

  <div class="calendar">
    {#each days as day}
      <div
        class="cell {selectedDate && day && formatDate(selectedDate) === formatDate(day) ? 'selected' : ''}"
        on:click={() => day && selectDay(day)}
        on:contextmenu={(e) => day && openDayContextMenu(e, day)}
      >
        {#if day}
          <div class="date">{day.getDate()}</div>

          <div class="events">
            {#each events.filter((e) => e.date === formatDate(day)) as event}
              <div
                class="event"
                style="background:{event.color}"
                on:contextmenu={(e) => openEventContextMenu(e, event)}
              >
                {event.title}
              </div>
            {/each}
          </div>
        {/if}
      </div>
    {/each}
  </div>
</div>

{#if contextMenu.visible}
  <div
    class="context-menu"
    style="top:{contextMenu.y}px; left:{contextMenu.x}px"
    on:click|stopPropagation={() => {}}
  >
    {#if contextMenu.kind === "day"}
      <div class="menu-item" on:click={addEvent}>➕ Add Event</div>
    {:else}
      <div class="menu-item" on:click={deleteEventFromMenu}>🗑 Sil</div>
    {/if}
  </div>
{/if}

<style>
.container {
  width: 100%;
  max-width: 900px;
  margin: auto;
  padding: 20px;
  color: white;
  font-family: sans-serif;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.header button {
  background: #2a2a2a;
  border: none;
  color: white;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
}

.header button:hover {
  background: #3a3a3a;
}

.weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  margin-bottom: 8px;
  text-align: center;
  opacity: 0.6;
}

.calendar {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 8px;
}

.cell {
  height: 110px;
  background: #1e1e1e;
  border-radius: 12px;
  padding: 6px;
  cursor: pointer;
  transition: 0.2s;
  position: relative;
}

.cell:hover {
  background: #2a2a2a;
  transform: scale(1.03);
}

.cell.selected {
  outline: 2px solid #4ade80;
}

.date {
  font-size: 14px;
  margin-bottom: 4px;
}

.events {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.event {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 6px;
  color: black;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: context-menu;
}

.context-menu {
  position: fixed;
  background: #2a2a2a;
  border-radius: 10px;
  padding: 6px;
  z-index: 1000;
  box-shadow: 0 10px 25px rgba(0,0,0,0.5);
}

.menu-item {
  padding: 8px 12px;
  cursor: pointer;
  border-radius: 6px;
}

.menu-item:hover {
  background: #3a3a3a;
}
</style>
