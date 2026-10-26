<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Card from '$lib/components/ui/Card.svelte';
  import Checkbox from '$lib/components/ui/Checkbox.svelte';
  import { i18n } from '$lib/i18n.svelte';
  

import Terminal from "lucide-svelte/icons/terminal";
import AlertCircle from "lucide-svelte/icons/alert-circle";
import User from "lucide-svelte/icons/user";
import Lock from "lucide-svelte/icons/lock";



  let username = $state('');
  let password = $state('');
  let remember = $state(false);
  let errorMsg = $state('');
  let loading = $state(false);
  let oidcSources = $state<any[]>([]);

  async function fetchSettings() {
    try {
      const res = await apiFetch('/api/v2/user/login/settings');
      if (res.ok) {
        const data = await res.json();
        oidcSources = data.oidc_sources || [];
      }
    } catch (e) {}
  }

  onMount(() => {
    fetchSettings();
  });

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
  <div class="background-decor">
    <div class="circle c1"></div>
    <div class="circle c2"></div>
    <div class="circle c3"></div>
  </div>

  <div class="login-container">
    <Card variant="highest" padding="none" class="login-card">
      <div class="login-layout">
        <div class="login-brand">
          <div class="brand-content">
            <Terminal size={64} class="logo-icon" />
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
              <AlertCircle size={20} />
              <p>{errorMsg}</p>
            </div>
          {/if}

          <form onsubmit={handleLogin} class="auth-form">
            <Input 
              id="username"
              label="Usuario o Email" 
              bind:value={username} 
              icon={User} 
              placeholder="tu_usuario" 
              required
            />
            
            <Input 
              id="password"
              label="Contraseña" 
              type="password" 
              bind:value={password} 
              icon={Lock} 
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

          {#if oidcSources.length > 0}
            <div class="oauth-divider">
              <span>O continuar con</span>
            </div>
            
            <div class="oauth-buttons">
              {#each oidcSources as source}
                <Button 
                  variant="outline" 
                  class="oidc-btn" 
                  href={`/api/v2/user/login/oidc/${source.id}`}
                >
                  {source.name}
                </Button>
              {/each}
            </div>
          {/if}

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
    position: relative;
    min-height: 100%;
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: hidden;
    background: transparent;
  }

  .background-decor {
    position: absolute;
    inset: 0;
    overflow: hidden;
    z-index: 0;
    filter: blur(100px);
    opacity: 0.35;
    pointer-events: none;
  }

  .circle {
    position: absolute;
    border-radius: 50%;
    animation: float 25s infinite alternate ease-in-out;
  }

  .c1 {
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, var(--primary-color), var(--status-info));
    top: -150px;
    right: -150px;
    animation-delay: 0s;
  }

  .c2 {
    width: 500px;
    height: 500px;
    background: radial-gradient(circle, var(--status-accent), var(--status-accent));
    bottom: -150px;
    left: -100px;
    animation-delay: -5s;
    animation-duration: 30s;
  }

  .c3 {
    width: 400px;
    height: 400px;
    background: radial-gradient(circle, var(--status-info), var(--status-success));
    top: 50%;
    left: -200px;
    transform: translateY(-50%);
    animation-delay: -10s;
    animation-duration: 35s;
  }

  @keyframes float {
    0% { transform: translate(0, 0) rotate(0deg); }
    100% { transform: translate(60px, 40px) rotate(10deg); }
  }

  .login-container {
    position: relative;
    z-index: 10;
    width: 100%;
    max-width: 900px;
    padding: 2rem;
    animation: slideUp 0.8s cubic-bezier(0.2, 0.8, 0.2, 1);
  }

  :global(.login-card) {
    position: relative;
    border-radius: 2.5rem !important;
    overflow: hidden !important;
    box-shadow: 
      0 32px 80px rgba(0, 0, 0, 0.4),
      inset 0 1px 1px rgba(255, 255, 255, 0.1),
      inset 0 -1px 1px rgba(0, 0, 0, 0.2) !important;
    border: 1px solid var(--md-sys-color-outline-variant) !important;
    background: var(--md-sys-color-surface-container-highest) !important;
  }

  .login-layout {
    display: grid;
    grid-template-columns: 1fr 1.2fr;
    min-height: 600px;
  }

  .login-brand {
    padding: 3.5rem;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    background: var(--md-sys-color-surface-container-low);
    border-right: 1px solid var(--md-sys-color-outline-variant);
  }

  .brand-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  :global(.logo-icon) {
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
  
  .oauth-divider {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin: 2rem 0;
    color: var(--text-muted);
    font-size: 0.85rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }
  
  .oauth-divider::before,
  .oauth-divider::after {
    content: "";
    flex: 1;
    height: 1px;
    background: var(--md-sys-color-outline-variant);
  }
  
  .oauth-buttons {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  
  :global(.oidc-btn) {
    height: 3.5rem !important;
    border-radius: 1.25rem !important;
    font-weight: 700 !important;
    border-color: var(--md-sys-color-outline-variant) !important;
  }
  
  :global(.oidc-btn:hover) {
    background: var(--md-sys-color-surface-container-high) !important;
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
