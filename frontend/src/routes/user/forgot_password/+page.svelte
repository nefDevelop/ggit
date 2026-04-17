<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';

  let email = $state('');
  let message = $state('');
  let errorMsg = $state('');
  let loading = $state(false);

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }

  async function handleForgot(e: Event) {
    e.preventDefault();
    loading = true;
    errorMsg = '';
    message = '';

    try {
      const res = await apiFetch('/api/v2/user/forgot_password', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'X-Csrf-Token': getCsrfToken()
        },
        body: JSON.stringify({ email })
      });

      const data = await res.json();

      if (res.ok && data.status === 'ok') {
        message = data.message || 'If the email exists, a reset link has been sent.';
      } else {
        errorMsg = data.message || 'Failed to send recovery email. Please try again.';
      }
    } catch (err) {
      errorMsg = 'A network error occurred. Please try again.';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Forgot Password - Gogs</title>
</svelte:head>

<div class="container auth-container">
  <div class="glass auth-window">
    <div class="auth-header">
      <h2>Forgot Password?</h2>
      <p>Enter your email to receive a recovery link</p>
    </div>

    {#if message}
      <div class="alert box-success">
        {message}
      </div>
    {/if}

    {#if errorMsg}
      <div class="alert box-error">
        {errorMsg}
      </div>
    {/if}

    {#if !message}
        <form onsubmit={handleForgot} class="auth-form">
        <div class="form-group">
            <label for="email">Email Address</label>
            <input 
            type="email" 
            id="email" 
            bind:value={email} 
            required 
            autocomplete="email"
            placeholder="your@email.com"
            />
        </div>

        <button type="submit" class="btn btn-primary" disabled={loading}>
            {loading ? 'Sending link...' : 'Send Recovery Link'}
        </button>
        </form>
    {/if}

    <div class="auth-footer">
      <p>Remembered your password? <a href="/user/login">Back to Sign In</a></p>
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

  .box-success {
    background-color: rgba(34, 197, 94, 0.1);
    border: 1px solid #22c55e;
    color: #86efac;
  }
</style>
