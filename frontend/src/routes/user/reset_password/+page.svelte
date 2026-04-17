<script lang="ts">
  import { page } from '$app/state';
  import { onMount } from 'svelte';

  let password = $state('');
  let retype = $state('');
  let code = $state('');
  let errorMsg = $state('');
  let successMsg = $state('');
  let loading = $state(false);

  onMount(() => {
    code = page.url.searchParams.get('code') || '';
    if (!code) {
      errorMsg = 'Invalid or missing reset code. Please request a new link.';
    }
  });

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }

  async function handleReset(e: Event) {
    e.preventDefault();
    if (password !== retype) {
        errorMsg = 'Passwords do not match.';
        return;
    }
    if (password.length < 6) {
        errorMsg = 'Password must be at least 6 characters.';
        return;
    }

    loading = true;
    errorMsg = '';
    successMsg = '';

    try {
      const res = await fetch('/api/v2/user/reset_password', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'X-Csrf-Token': getCsrfToken()
        },
        body: JSON.stringify({ code, password })
      });

      const data = await res.json();

      if (res.ok && data.status === 'ok') {
        successMsg = 'Password has been successfully reset. You can now sign in.';
      } else {
        errorMsg = data.message || 'Failed to reset password. The link might be expired.';
      }
    } catch (err) {
      errorMsg = 'A network error occurred. Please try again.';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Reset Password - Gogs</title>
</svelte:head>

<div class="container auth-container">
  <div class="glass auth-window">
    <div class="auth-header">
      <h2>Reset Password</h2>
      <p>Enter your new password below</p>
    </div>

    {#if successMsg}
      <div class="alert box-success">
        {successMsg}
      </div>
      <div class="auth-footer">
        <a href="/user/login" class="btn btn-primary">Go to Sign In</a>
      </div>
    {:else}
        {#if errorMsg}
            <div class="alert box-error">
                {errorMsg}
            </div>
        {/if}

        {#if code}
            <form onsubmit={handleReset} class="auth-form">
                <div class="form-group">
                    <label for="password">New Password</label>
                    <input 
                    type="password" 
                    id="password" 
                    bind:value={password} 
                    required 
                    autocomplete="new-password"
                    placeholder="••••••••"
                    />
                </div>

                <div class="form-group">
                    <label for="retype">Confirm Password</label>
                    <input 
                    type="password" 
                    id="retype" 
                    bind:value={retype} 
                    required 
                    autocomplete="new-password"
                    placeholder="••••••••"
                    />
                </div>

                <button type="submit" class="btn btn-primary" disabled={loading}>
                    {loading ? 'Resetting...' : 'Reset Password'}
                </button>
            </form>
        {/if}
    {/if}

    <div class="auth-footer">
      <p>Back to <a href="/user/login">Sign In</a></p>
    </div>
  </div>
</div>

<style>
  /* Reuse same styles as forgot_password */
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
    display: inline-block;
    text-align: center;
    text-decoration: none;
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
