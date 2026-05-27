(() => {
    // Prevent duplicate initialization
    if (window.alertNew || window.confirmNew) {
      return;
    }
  
    // =========================
    // Styles
    // =========================
  
    const style = document.createElement("style");
  
    style.textContent = `
      #__native_modal_overlay__ {
        position: fixed;
        inset: 0;
        background: rgba(0,0,0,0.6);
        z-index: 2147483647;
        display: none;
        align-items: center;
        justify-content: center;
        backdrop-filter: blur(2px);
      }
  
      #__native_modal__ {
        width: min(800px, 90vw);
        max-height: 85vh;
        background: #1b1b1b;
        color: white;
        border-radius: 14px;
        overflow: hidden;
        display: flex;
        flex-direction: column;
        box-shadow: 0 10px 40px rgba(0,0,0,0.5);
        font-family: system-ui, sans-serif;
      }
  
      #__native_modal_header__ {
        padding: 14px 18px;
        background: #2a2a2a;
        border-bottom: 1px solid #3a3a3a;
        font-size: 15px;
        font-weight: bold;
      }
  
      #__native_modal_content__ {
        padding: 18px;
        overflow: auto;
        white-space: pre-wrap;
        word-break: break-word;
        line-height: 1.5;
        font-family: monospace;
        font-size: 13px;
        flex: 1;
      }
  
      #__native_modal_actions__ {
        display: flex;
        justify-content: flex-end;
        gap: 10px;
        padding: 16px;
        border-top: 1px solid #3a3a3a;
        background: #222;
      }
  
      .__native_modal_btn__ {
        border: none;
        border-radius: 8px;
        padding: 10px 16px;
        cursor: pointer;
        font-size: 14px;
        min-width: 90px;
      }
  
      #__native_modal_ok__ {
        background: #0a84ff;
        color: white;
      }
  
      #__native_modal_cancel__ {
        background: #444;
        color: white;
      }
  
      #__native_modal_close__ {
        background: #444;
        color: white;
      }
  
      .__native_modal_btn__:hover {
        filter: brightness(1.1);
      }
    `;
  
    document.head.appendChild(style);
  
    // =========================
    // DOM
    // =========================
  
    const overlay = document.createElement("div");
  
    overlay.id = "__native_modal_overlay__";
  
    overlay.innerHTML = `
      <div id="__native_modal__">
        <div id="__native_modal_header__"></div>
  
        <div id="__native_modal_content__"></div>
  
        <div id="__native_modal_actions__">
          <button
            class="__native_modal_btn__"
            id="__native_modal_cancel__"
          >
            Cancel
          </button>
  
          <button
            class="__native_modal_btn__"
            id="__native_modal_close__"
          >
            Close
          </button>
  
          <button
            class="__native_modal_btn__"
            id="__native_modal_ok__"
          >
            OK
          </button>
        </div>
      </div>
    `;
  
    document.body.appendChild(overlay);
  
    const header =
      document.getElementById("__native_modal_header__");
  
    const content =
      document.getElementById("__native_modal_content__");
  
    const actions =
      document.getElementById("__native_modal_actions__");
  
    const okBtn =
      document.getElementById("__native_modal_ok__");
  
    const cancelBtn =
      document.getElementById("__native_modal_cancel__");
  
    const closeBtn =
      document.getElementById("__native_modal_close__");
  
    let resolver = null;
    let currentMode = "alert";
  
    // =========================
    // Helpers
    // =========================
  
    function setContent(value) {
      content.textContent =
        typeof value === "string"
          ? value
          : JSON.stringify(value, null, 2);
    }
  
    function close(result = false) {
      overlay.style.display = "none";
  
      if (resolver) {
        resolver(result);
        resolver = null;
      }
    }
  
    function open({
      mode = "alert",
      title = "Message",
      message = "",
    }) {
      currentMode = mode;
  
      header.textContent = title;
  
      setContent(message);
  
      overlay.style.display = "flex";
  
      // Alert mode
      if (mode === "alert") {
        okBtn.style.display = "none";
        cancelBtn.style.display = "none";
        closeBtn.style.display = "inline-block";
  
        return new Promise((resolve) => {
          resolver = resolve;
        });
      }
  
      // Confirm mode
      okBtn.style.display = "inline-block";
      cancelBtn.style.display = "inline-block";
      closeBtn.style.display = "none";
  
      return new Promise((resolve) => {
        resolver = resolve;
      });
    }
  
    // =========================
    // Events
    // =========================
  
    okBtn.addEventListener("click", () => {
      close(true);
    });
  
    cancelBtn.addEventListener("click", () => {
      close(false);
    });
  
    closeBtn.addEventListener("click", () => {
      close(true);
    });
  
    overlay.addEventListener("click", (e) => {
      if (e.target !== overlay) {
        return;
      }
  
      if (currentMode === "confirm") {
        close(false);
      } else {
        close(true);
      }
    });
  
    document.addEventListener("keydown", (e) => {
      if (overlay.style.display !== "flex") {
        return;
      }
  
      if (e.key === "Escape") {
        if (currentMode === "confirm") {
          close(false);
        } else {
          close(true);
        }
      }
  
      if (e.key === "Enter") {
        close(true);
      }
    });
  
    // =========================
    // Public API
    // =========================
  
    window.alertNew = function (
      message,
      title = "Message"
    ) {
      return open({
        mode: "alert",
        title,
        message,
      });
    };
  
    window.confirmNew = function (
      message,
      title = "Confirm"
    ) {
      return open({
        mode: "confirm",
        title,
        message,
      });
    };
  
    // =========================
    // Console helper
    // =========================
  
    console.log(`
  alertNew("Hello world")
  
  confirmNew("Delete file?")
    .then(result => console.log(result))
  `);
  })();