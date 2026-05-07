<script lang="ts">
  import { onMount } from "svelte";
  import { apiFetch } from "$lib/api";
  import { i18n } from "$lib/i18n.svelte";
  import { dialog } from "$lib/dialog.svelte";
  import { KeyRound, CheckCircle2, AlertCircle, SmartphoneNfc, ShieldCheck, ShieldAlert, ShieldPlus, Key, ShieldOff } from "lucide-svelte";
  import Card from "./ui/Card.svelte";
  import Button from "./ui/Button.svelte";

  let oldPassword = $state("");
  let newPassword = $state("");
  let repeatPassword = $state("");
  let saving = $state(false);
  let message = $state("");
  let isError = $state(false);

  // 2FA State
  let twoFactorEnabled = $state(false);
  let enrollment = $state(null as any);
  let passcode = $state("");
  let recoveryCodes = $state([] as string[]);
  let showRecovery = $state(false);

  async function updatePassword() {
    if (newPassword !== repeatPassword) {
      message = "Las contraseñas no coinciden.";
      isError = true;
      return;
    }
    saving = true;
    message = "";
    isError = false;
    try {
      const res = await apiFetch("/api/v2/user/settings/password", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          old_password: oldPassword,
          password: newPassword,
          retype: repeatPassword,
        }),
      });
      if (res.ok) {
        message = "¡Contraseña actualizada correctamente!";
        oldPassword = "";
        newPassword = "";
        repeatPassword = "";
      } else {
        const body = await res.json();
        message = body.message || "Error al actualizar la contraseña.";
        isError = true;
      }
    } catch (e) {
      message = "Error de red.";
      isError = true;
    } finally {
      saving = false;
    }
  }

  async function fetch2FAStatus() {
    try {
      const res = await apiFetch("/api/v2/user/settings/security/two_factor");
      if (res.ok) {
        const data = await res.json();
        twoFactorEnabled = data.enabled;
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function initiate2FA() {
    try {
      const res = await apiFetch(
        "/api/v2/user/settings/security/two_factor/enable",
        { method: "POST" },
      );
      if (res.ok) enrollment = await res.json();
    } catch (e) {
      console.error(e);
    }
  }

  async function confirm2FA() {
    try {
      const res = await apiFetch(
        "/api/v2/user/settings/security/two_factor/confirm",
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ passcode }),
        },
      );
      if (res.ok) {
        const data = await res.json();
        recoveryCodes = data.recovery_codes;
        twoFactorEnabled = true;
        enrollment = null;
        showRecovery = true;
      } else {
        alert("Código inválido");
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function disable2FA() {
    if (
      !(await dialog.confirm(
        "¿Seguro que quieres deshabilitar el 2FA? Tu cuenta será menos segura.",
      ))
    )
      return;
    try {
      const res = await apiFetch("/api/v2/user/settings/security/two_factor", {
        method: "DELETE",
      });
      if (res.ok) {
        twoFactorEnabled = false;
        showRecovery = false;
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetch2FAStatus);
</script>

<div class="settings-security-container">

  <div class="modular-stack">
    <!-- PASSWORD SECTION -->
    <div class="config-module">
      <div class="module-tag">
        <KeyRound size={20} />
        Cambiar Contraseña
      </div>
      <div class="island form-box">
        <form
          onsubmit={(e) => {
            e.preventDefault();
            updatePassword();
          }}
          class="security-form"
        >
          <div class="form-grid">
            <div class="field">
              <label for="oldPassword">Contraseña Actual</label>
              <input
                id="oldPassword"
                type="password"
                bind:value={oldPassword}
                placeholder="••••••••"
                required
              />
            </div>
            <div class="field">
              <label for="newPassword">Nueva Contraseña</label>
              <input
                id="newPassword"
                type="password"
                bind:value={newPassword}
                placeholder="Mínimo 8 caracteres"
                required
              />
            </div>
            <div class="field">
              <label for="repeatPassword">Confirmar Nueva Contraseña</label>
              <input
                id="repeatPassword"
                type="password"
                bind:value={repeatPassword}
                placeholder="••••••••"
                required
              />
            </div>
          </div>

          <div class="form-footer">
            {#if message}
              <div class="status-alert" class:is-error={isError}>
                {#if isError}
                  <AlertCircle size={20} />
                {:else}
                  <CheckCircle2 size={20} />
                {/if}
                {message}
              </div>
            {/if}
            <Button
              type="submit"
              variant="primary"
              disabled={saving || !oldPassword}
            >
              {saving ? "Actualizando..." : "Actualizar Contraseña"}
            </Button>
          </div>
        </form>
      </div>
    </div>

    <!-- 2FA SECTION -->
    <div class="config-module">
      <div class="module-tag">
        <SmartphoneNfc size={20} />
        Doble Factor (2FA)
      </div>
      <div class="island enrollment-box" class:is-protected={twoFactorEnabled}>
        <div class="enrollment-header">
          <div class="status-indicator">
            {#if twoFactorEnabled}
              <ShieldCheck size={48} />
            {:else}
              <ShieldAlert size={48} />
            {/if}
            <div class="status-text">
              <strong
                >{twoFactorEnabled
                  ? "Cuenta Protegida"
                  : "Protección Deshabilitada"}</strong
              >
              <span
                >{twoFactorEnabled
                  ? "El acceso requiere un código temporal adicional."
                  : "Añade una capa extra de seguridad a tu cuenta."}</span
              >
            </div>
          </div>
          {#if !twoFactorEnabled && !enrollment}
            <Button
              variant="primary"
              icon={ShieldPlus}
              onclick={initiate2FA}>Configurar 2FA</Button
            >
          {/if}
        </div>

        {#if showRecovery}
          <div class="recovery-view">
            <h4>¡Éxito! 2FA Habilitado</h4>
            <p>
              Guarda estos códigos de recuperación en un lugar seguro. Los
              necesitarás si pierdes acceso a tu dispositivo.
            </p>
            <div class="recovery-grid">
              {#each recoveryCodes as code}
                <code>{code}</code>
              {/each}
            </div>
            <Button variant="outline" onclick={() => (showRecovery = false)}
              >He guardado los códigos</Button
            >
          </div>
        {:else if enrollment}
          <div class="enrollment-flow">
            <div class="steps-row">
              <div class="step-card tonal">
                <span class="step-num">1</span>
                <p>
                  Escanea el código QR con una app como Google Authenticator o
                  Authy.
                </p>
                <div class="qr-wrap">
                  <img src={enrollment.qr_code} alt="QR" />
                </div>
                <div class="manual-code">
                  <small>O introduce el secreto manualmente:</small>
                  <code>{enrollment.secret}</code>
                </div>
              </div>
              <div class="step-card tonal">
                <span class="step-num">2</span>
                <p>
                  Introduce el código de 6 dígitos que aparece en tu aplicación
                  móvil.
                </p>
                <div class="verify-box">
                  <input
                    type="text"
                    maxlength="6"
                    placeholder="000000"
                    bind:value={passcode}
                  />
                  <Button
                    variant="primary"
                    onclick={confirm2FA}
                    disabled={passcode.length < 6}>Verificar y Activar</Button
                  >
                </div>
                <Button variant="ghost" onclick={() => (enrollment = null)}
                  >Cancelar Proceso</Button
                >
              </div>
            </div>
          </div>
        {:else if twoFactorEnabled}
          <div class="manage-2fa">
            <div class="actions-row">
              <Button
                variant="outline"
                icon={Key}
                onclick={async () => {
                  const res = await apiFetch(
                    "/api/v2/user/settings/security/two_factor/recovery_codes",
                  );
                  if (res.ok) {
                    const data = await res.json();
                    recoveryCodes = data.recovery_codes;
                    showRecovery = true;
                  }
                }}>Ver Códigos de Recuperación</Button
              >
              <Button
                variant="ghost"
                icon={ShieldOff}
                class="danger-btn"
                onclick={disable2FA}>Desactivar 2FA</Button
              >
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>

<style>


  .modular-stack {
    display: flex;
    flex-direction: column;
    gap: 4rem;
  }
  .config-module {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .module-tag {
    font-size: 0.8rem;
    font-weight: 950;
    text-transform: uppercase;
    color: var(--primary-color);
    letter-spacing: 0.15em;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    border-bottom: 2px solid var(--md-sys-color-outline-variant);
    padding-bottom: 0.75rem;
    width: fit-content;
    opacity: 0.8;
  }
  .module-tag :global(svg) {
    font-size: 1.25rem;
  }

  .island {
    border-radius: 2.5rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    background: rgba(255, 255, 255, 0.01);
  }
  .form-box {
    padding: 3rem;
  }
  @media (max-width: 600px) {
    .form-box {
      padding: 1.5rem;
      border-radius: 1.5rem;
    }
  }

  .security-form {
    display: flex;
    flex-direction: column;
    gap: 2.5rem;
  }
  .form-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(100%, 250px), 1fr));
    gap: 2rem;
  }
  .field {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
  .field label {
    font-size: 0.85rem;
    font-weight: 800;
    color: var(--text-color);
    margin-left: 0.5rem;
  }
  .field input {
    background: var(--surface-variant);
    border: 2px solid transparent;
    padding: 1rem 1.5rem;
    border-radius: 1.25rem;
    color: var(--text-color);
    font-size: 1.1rem;
    transition: all 0.3s;
    width: 100%;
    box-sizing: border-box;
  }
  .field input:focus {
    border-color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.05);
  }

  .form-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 2rem;
    border-top: 1px solid var(--md-sys-color-outline-variant);
    gap: 1.5rem;
  }
  @media (max-width: 800px) {
    .form-footer {
      flex-direction: column;
      align-items: stretch;
    }
    .status-alert {
      order: 1;
    }
  }
  .status-alert {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-weight: 700;
    color: #10b981;
  }
  .status-alert.is-error {
    color: var(--error);
  }

  /* 2FA ENROLLMENT */
  .enrollment-box {
    padding: 3rem;
  }
  @media (max-width: 600px) {
    .enrollment-box {
      padding: 1.5rem;
      border-radius: 1.5rem;
    }
  }
  .enrollment-box.is-protected {
    background: rgba(16, 185, 129, 0.03);
    border-color: rgba(16, 185, 129, 0.2);
  }

  .enrollment-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 2rem;
  }
  .status-indicator {
    display: flex;
    align-items: center;
    gap: 1.5rem;
  }
  .status-indicator span {
    font-size: 3rem;
    color: var(--text-muted);
  }
  .is-protected .status-indicator span {
    color: #10b981;
  }
  .status-text {
    display: flex;
    flex-direction: column;
  }
  .status-text strong {
    font-size: 1.5rem;
    font-weight: 900;
  }
  .status-text span {
    color: var(--text-muted);
    font-size: 1rem;
  }

  .recovery-view {
    margin-top: 2.5rem;
    padding: 2.5rem;
    background: #08080c;
    border-radius: 2rem;
    border: 1px solid #1a1a24;
    text-align: center;
  }
  @media (max-width: 600px) {
    .recovery-view {
      padding: 1.5rem;
    }
  }
  .recovery-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 0.75rem;
    margin: 2rem 0;
  }
  @media (max-width: 500px) {
    .recovery-grid {
      grid-template-columns: 1fr;
    }
  }
  .recovery-grid code {
    background: #000;
    padding: 1rem;
    border-radius: 1rem;
    color: var(--primary-color);
    font-weight: 900;
    font-family: "JetBrains Mono", monospace;
    font-size: 1.1rem;
    letter-spacing: 0.1em;
    border: 1px solid #222;
  }

  .enrollment-flow {
    margin-top: 3rem;
  }
  .steps-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 2rem;
  }
  .step-card {
    padding: 2rem;
    border-radius: 2rem;
    position: relative;
  }
  @media (max-width: 600px) {
    .step-card {
      padding: 1.5rem;
    }
  }
  .step-num {
    position: absolute;
    top: 1.5rem;
    right: 1.5rem;
    width: 32px;
    height: 32px;
    background: var(--primary-color);
    color: white;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 900;
  }
  .qr-wrap {
    background: white;
    padding: 1rem;
    border-radius: 1.5rem;
    width: fit-content;
    margin: 1.5rem auto;
  }
  .qr-wrap img {
    width: 180px;
    height: 180px;
  }
  .manual-code {
    text-align: center;
  }
  .manual-code code {
    display: block;
    margin-top: 0.5rem;
    background: #000;
    padding: 0.75rem;
    border-radius: 0.75rem;
    color: var(--primary-color);
    font-weight: 800;
    word-break: break-all;
  }

  .verify-box {
    display: flex;
    flex-direction: column;
    gap: 1rem;
    margin: 2rem 0;
  }
  .verify-box input {
    text-align: center;
    font-size: 2rem;
    font-weight: 950;
    letter-spacing: 0.3em;
    background: #000;
    border: 2px solid var(--primary-color);
    padding: 1rem;
    border-radius: 1.25rem;
    color: white;
    width: 100%;
    box-sizing: border-box;
  }
  @media (max-width: 600px) {
    .verify-box input {
      font-size: 1.5rem;
      letter-spacing: 0.2em;
    }
  }

  .manage-2fa {
    margin-top: 2.5rem;
    border-top: 1px solid var(--md-sys-color-outline-variant);
    padding-top: 2.5rem;
  }
  .actions-row {
    display: flex;
    gap: 1rem;
    flex-wrap: wrap;
  }
  :global(.danger-btn:hover) {
    color: var(--error) !important;
    background: rgba(var(--error-rgb), 0.1) !important;
  }

  @media (max-width: 900px) {
    .steps-row {
      grid-template-columns: 1fr;
    }
    .enrollment-header {
      flex-direction: column;
      align-items: stretch;
      text-align: center;
    }
    .status-indicator {
      flex-direction: column;
    }
  }
</style>
