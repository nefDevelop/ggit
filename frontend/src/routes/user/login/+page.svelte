<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Card from '$lib/components/ui/Card.svelte';
  import Checkbox from '$lib/components/ui/Checkbox.svelte';
  import { i18n } from '$lib/i18n.svelte';

  let username = $state('');
  let password = $state('');
  let remember = $state(false);
  let errorMsg = $state('');
  let loading = $state(false);

  async function handleLogin(e: Event) {
    e.preventDefault();
    loading = true;
    errorMsg = '';

    try {
      const res = await apiFetch('/api/v2/user/login', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username,
          password,
          remember,
          login_source: 0 // Default local
        })
      });

      const data = await res.json();

      if (res.ok && data.status === 'ok') {
        authState.user = data.user;
        authState.isAuthenticated = true;
        window.location.href = '/dashboard'; // Dashboard
      } else {
        errorMsg = data.message || 'Login failed. Please check your credentials.';
      }
    } catch (err) {
      errorMsg = 'A network error occurred. Please try again.';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Sign In - Ggit</title>
</svelte:head>

<div class="login-page">
  <div class="login-container">
    <Card variant="highest" padding="none" class="login-card glass">
      <div class="background-decor">
        <div class="circle c1"></div>
        <div class="circle c2"></div>
        <div class="circle c3"></div>
      </div>

      <div class="login-layout">
        <div class="login-brand tonal">
          <div class="brand-content">
            <span class="material-symbols-rounded logo-icon">terminal</span>
            <h1>Ggit</h1>
            <p>El servicio de Git más ligero y elegante, ahora rediseñado para ti.</p>
          </div>
          <div class="brand-footer">
            <span>v2026.4.0</span>
          </div>
        </div>

        <div class="login-form-area">
          <header class="form-header">
            <h2>Bienvenido</h2>
            <p>Inicia sesión para continuar</p>
          </header>

          {#if errorMsg}
            <div class="error-banner animate-in">
              <span class="material-symbols-rounded">error</span>
              <p>{errorMsg}</p>
            </div>
          {/if}

          <form onsubmit={handleLogin} class="auth-form">
            <Input 
              label="Usuario o Email" 
              bind:value={username} 
              icon="person" 
              placeholder="tu_usuario" 
              required
            />
            
            <Input 
              label="Contraseña" 
              type="password" 
              bind:value={password} 
              icon="lock" 
              placeholder="••••••••" 
              required
            />

            <div class="form-options">
              <Checkbox label="Recordarme" bind:checked={remember} simple />
              <a href="/user/forgot_password" class="forgot-link">¿Olvidaste tu contraseña?</a>
            </div>

            <Button type="submit" variant="primary" class="login-btn" loading={loading}>
              Iniciar Sesión
            </Button>
          </form>

          <footer class="form-footer">
            <p>¿No tienes una cuenta? <a href="/user/sign_up">Regístrate ahora</a></p>
          </footer>
        </div>
      </div>
    </Card>
  </div>
</div>

<style>
  .login-page {
    position: fixed;
    inset: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    background: var(--bg-color);
    overflow: hidden;
  }

  .background-decor {
    position: absolute;
    inset: 0;
    overflow: hidden;
    z-index: 0;
    filter: blur(80px);
    opacity: 0.5;
    pointer-events: none;
  }

  .circle {
    position: absolute;
    border-radius: 50%;
  }

  .c1 {
    width: 400px;
    height: 400px;
    background: radial-gradient(circle, var(--primary-color), #3b82f6);
    top: -150px;
    right: -100px;
  }

  .c2 {
    width: 350px;
    height: 350px;
    background: radial-gradient(circle, #8b5cf6, #d946ef);
    bottom: -120px;
    left: -80px;
  }

  .c3 {
    width: 250px;
    height: 250px;
    background: radial-gradient(circle, #06b6d4, #10b981);
    top: 50%;
    left: -100px;
    transform: translateY(-50%);
  }

  .login-container {
    position: relative;
    z-index: 1;
    width: 100%;
    max-width: 900px;
    padding: 2rem;
    animation: slideUp 0.6s cubic-bezier(0.2, 0.8, 0.2, 1);
  }

  :global(.login-card) {
    position: relative;
    border-radius: 2.5rem !important;
    overflow: hidden !important;
    box-shadow: 0 24px 48px rgba(0, 0, 0, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .login-layout {
    display: grid;
    grid-template-columns: 1fr 1.2fr;
    min-height: 600px;
    position: relative;
    z-index: 1;
  }

  .login-brand {
    padding: 3rem;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    position: relative;
    overflow: hidden;
    background: rgba(var(--primary-rgb), 0.05);
    backdrop-filter: blur(10px);
  }

  .brand-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .logo-icon {
    font-size: 4rem;
    color: var(--primary-color);
  }

  .login-brand h1 {
    font-size: 3.5rem;
    font-weight: 950;
    margin: 0;
    letter-spacing: -0.05em;
    background: linear-gradient(to bottom, var(--text-color), var(--primary-color));
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .login-brand p {
    font-size: 1.1rem;
    line-height: 1.6;
    color: var(--text-muted);
    max-width: 300px;
  }

  .brand-footer {
    font-size: 0.8rem;
    font-weight: 800;
    color: var(--text-muted);
    opacity: 0.5;
  }

  .login-form-area {
    padding: 3.5rem;
    display: flex;
    flex-direction: column;
    justify-content: center;
    background: var(--md-sys-color-surface-container-highest);
  }

  .form-header {
    margin-bottom: 2.5rem;
  }

  .form-header h2 {
    font-size: 2.25rem;
    font-weight: 900;
    margin: 0;
    letter-spacing: -0.02em;
  }

  .form-header p {
    color: var(--text-muted);
    margin-top: 0.5rem;
    font-size: 1.1rem;
  }

  .auth-form {
    display: flex;
    flex-direction: column;
    gap: 1.75rem;
  }

  .form-options {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-top: 0.5rem;
  }

  .forgot-link {
    font-size: 0.9rem;
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 600;
  }

  .forgot-link:hover {
    text-decoration: underline;
  }

  :global(.login-btn) {
    margin-top: 1rem;
    height: 3.5rem !important;
    font-size: 1.1rem !important;
    font-weight: 800 !important;
    border-radius: 1.25rem !important;
  }

  .error-banner {
    background: rgba(var(--error-rgb), 0.1);
    color: var(--error);
    padding: 1rem 1.25rem;
    border-radius: 1rem;
    margin-bottom: 2rem;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-weight: 600;
    font-size: 0.95rem;
  }

  .error-banner p { margin: 0; }

  .form-footer {
    margin-top: 2.5rem;
    text-align: center;
    font-size: 0.95rem;
    color: var(--text-muted);
  }

  .form-footer a {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 700;
  }

  .form-footer a:hover {
    text-decoration: underline;
  }

  @keyframes slideUp {
    from { opacity: 0; transform: translateY(30px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .animate-in {
    animation: fadeIn 0.3s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: scale(0.95); }
    to { opacity: 1; transform: scale(1); }
  }

  @media (max-width: 850px) {
    .login-layout {
      grid-template-columns: 1fr;
    }
    .login-brand {
      display: none;
    }
    .login-container {
      max-width: 500px;
    }
    .login-form-area {
      padding: 2.5rem;
    }
  }
</style>

