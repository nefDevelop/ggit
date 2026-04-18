<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';

  let username = $state('');
  let password = $state('');
  let remember = $state(false);
  let errorMsg = $state('');
  let loading = $state(false);

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }

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
  <title>Sign In - Gogs</title>
</svelte:head>

<div class="container auth-container">
  <div class="glass auth-window">
    <div class="auth-header">
      <h2>Welcome Back</h2>
      <p>Sign in to your account</p>
    </div>

    {#if errorMsg}
      <div class="alert box-error">
        {errorMsg}
      </div>
    {/if}

    <form onsubmit={handleLogin} class="auth-form">
      <div class="form-group">
        <label for="username">Username or Email</label>
        <input 
          type="text" 
          id="username" 
          bind:value={username} 
          required 
          autocomplete="username"
          placeholder="your_username"
        />
      </div>

      <div class="form-group">
        <label for="password">Password</label>
        <input 
          type="password" 
          id="password" 
          bind:value={password} 
          required 
          autocomplete="current-password"
          placeholder="••••••••"
        />
      </div>

      <div class="form-actions">
        <label class="remember-me">
          <input type="checkbox" bind:checked={remember} />
          <span>Remember me</span>
        </label>
        <a href="/user/forgot_password" class="forgot-link">Forgot password?</a>
      </div>

      <button type="submit" class="btn btn-primary" disabled={loading}>
        {loading ? 'Signing in...' : 'Sign In'}
      </button>
    </form>

    <div class="auth-footer">
      <p>Don't have an account? <a href="/user/sign_up">Sign up now</a></p>
    </div>
  </div>
</div>

<style>
  .auth-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: calc(100vh - 100px);
  }

  .auth-window {
    width: 100%;
    max-width: 420px;
    padding: 3rem 2.5rem;
  }

  .auth-header {
    text-align: center;
    margin-bottom: 2rem;
  }

  .auth-header h2 {
    font-size: 1.8rem;
    margin-bottom: 0.5rem;
  }

  .auth-header p {
    color: var(--text-muted);
  }

  .auth-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .form-group label {
    font-size: 0.9rem;
    font-weight: 500;
  }

  .form-group input {
    background-color: rgba(15, 23, 42, 0.5);
    border: 1px solid var(--border-color);
    color: var(--text-color);
    padding: 0.8rem 1rem;
    border-radius: var(--radius);
    outline: none;
    transition: border-color 0.2s ease;
  }

  .form-group input:focus {
    border-color: var(--primary-color);
  }

  .btn-primary {
    width: 100%;
    margin-top: 1rem;
    padding: 1rem;
    font-size: 1.1rem;
  }

  .btn-primary:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  .form-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.9rem;
  }

  .remember-me {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
  }

  .forgot-link {
    color: var(--primary-color);
    text-decoration: none;
  }

  .forgot-link:hover {
    text-decoration: underline;
  }

  .auth-footer {
    margin-top: 2rem;
    text-align: center;
    font-size: 0.9rem;
    color: var(--text-muted);
  }

  .auth-footer a {
    color: var(--primary-color);
    text-decoration: none;
    font-weight: 500;
  }

  .auth-footer a:hover {
    text-decoration: underline;
  }

  .alert {
    padding: 1rem;
    border-radius: var(--radius);
    margin-bottom: 1.5rem;
    font-size: 0.9rem;
  }

  .box-error {
    background-color: rgba(239, 68, 68, 0.1);
    border: 1px solid #ef4444;
    color: #fca5a5;
  }
</style>
