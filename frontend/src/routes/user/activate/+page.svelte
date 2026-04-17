<script lang="ts">
  import { page } from '$app/state';
  import { onMount } from 'svelte';
  import { authState } from '$lib/auth.svelte';

  let code = $state('');
  let email = $state('');
  let loading = $state(false);
  let errorMsg = $state('');
  let successMsg = $state('');
  let isActivated = $state(false);

  onMount(async () => {
    code = page.url.searchParams.get('code') || '';
    email = page.url.searchParams.get('email') || '';

    if (code) {
      handleActivate();
    }
  });

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }

  async function handleActivate() {
    loading = true;
    errorMsg = '';
    
    try {
      const res = await fetch('/api/v2/user/activate', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'X-Csrf-Token': getCsrfToken()
        },
        body: JSON.stringify({ code, email })
      });

      const data = await res.json();

      if (res.ok && data.status === 'ok') {
        isActivated = true;
        successMsg = 'Your account has been successfully activated!';
      } else {
        errorMsg = data.message || 'Activation failed. The link might be expired.';
      }
    } catch (err) {
      errorMsg = 'A network error occurred. Please try again.';
    } finally {
      loading = false;
    }
  }

  async function resendEmail() {
    loading = true;
    errorMsg = '';
    successMsg = '';

    try {
       // Note: There is no direct "resend" v2 API yet, but typically we can use forgot_password 
       // or a specific activate/resend endpoint. 
       // For now, let's assume the activate endpoint handles empty code as resend if email is provided.
       // Actually, I'll check v2/auth.go again.
       // Looking at internal/route/api/v2/auth.go: Activate doesn't handle resend.
       // But I'll implement it if needed.
       
       errorMsg = "Resend functionality is currently being modernized. Please contact support.";
    } catch (err) {
      errorMsg = 'Failed to resend email.';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Account Activation - Gogs</title>
</svelte:head>

<div class="container auth-container">
  <div class="glass auth-window">
    <div class="auth-header">
      <h2>Account Activation</h2>
      {#if !isActivated}
        <p>Verifying your account details...</p>
      {:else}
        <p>Registration complete!</p>
      {/if}
    </div>

    {#if loading}
      <div class="loading-state">
        <div class="spinner"></div>
        <p>Processing...</p>
      </div>
    {/if}

    {#if successMsg}
      <div class="alert box-success">
        {successMsg}
      </div>
      <div class="auth-footer">
        <a href="/dashboard" class="btn btn-primary">Go to Dashboard</a>
      </div>
    {:else if errorMsg}
      <div class="alert box-error">
        {errorMsg}
      </div>
      {#if email}
        <button onclick={resendEmail} class="btn btn-secondary" disabled={loading}>
            Resend Activation Email
        </button>
      {/if}
    {/if}

    {#if !code && !loading && !successMsg && !errorMsg}
        <div class="info-text">
            <p>Please check your email for the activation link.</p>
            {#if email}
                <p>Sent to: <strong>{email}</strong></p>
                <button onclick={resendEmail} class="btn btn-secondary" style="margin-top: 1rem;">
                    Resend Email
                </button>
            {/if}
        </div>
    {/if}

    <div class="auth-footer">
      <p>Need help? <a href="/user/login">Contact Support</a></p>
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

  .auth-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .btn-primary, .btn-secondary {
    width: 100%;
    margin-top: 1rem;
    padding: 1rem;
    font-size: 1.1rem;
    display: inline-block;
    text-align: center;
    text-decoration: none;
    border-radius: var(--radius);
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .btn-primary {
    background-color: var(--primary-color);
    color: white;
    border: none;
  }

  .btn-secondary {
    background-color: rgba(255, 255, 255, 0.05);
    border: 1px solid var(--border-color);
    color: var(--text-color);
  }

  .btn-secondary:hover {
    background-color: rgba(255, 255, 255, 0.1);
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

  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    margin: 2rem 0;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 4px solid rgba(255, 255, 255, 0.1);
    border-left-color: var(--primary-color);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .info-text {
      text-align: center;
      color: var(--text-muted);
  }
</style>
