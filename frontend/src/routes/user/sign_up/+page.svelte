<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';

  let username = $state('');
  let email = $state('');
  let password = $state('');
  let retype = $state('');
  let errorMsg = $state('');
  let successMsg = $state('');
  let loading = $state(false);

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }

  async function handleSignUp(e: Event) {
    e.preventDefault();
    if (password !== retype) {
      errorMsg = 'Passwords do not match.';
      return;
    }

    loading = true;
    errorMsg = '';
    successMsg = '';

    try {
      const res = await apiFetch('/api/v2/user/register', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'X-Csrf-Token': getCsrfToken()
        },
        body: JSON.stringify({
          username,
          email,
          password
        })
      });

      const data = await res.json();

      if (res.ok && data.status === 'ok') {
        if (data.requires_email_activation) {
          successMsg = 'Account created! Please check your email to activate your account.';
        } else {
          // Gogs doesn't auto-login after sign up, it redirects to login (or we can just redirect)
          successMsg = 'Account created successfully! Redirecting to login...';
          setTimeout(() => {
            window.location.href = '/user/login';
          }, 2000);
        }
      } else {
        errorMsg = data.message || 'Registration failed.';
      }
    } catch (err) {
      errorMsg = 'A network error occurred. Please try again.';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Sign Up - Gogs</title>
</svelte:head>

<div class="container auth-container">
  <div class="glass auth-window">
    <div class="auth-header">
      <h2>Create Account</h2>
      <p>Join the community and start collaborating</p>
    </div>

    {#if errorMsg}
      <div class="alert box-error">
        {errorMsg}
      </div>
    {/if}

    {#if successMsg}
      <div class="alert box-success">
        {successMsg}
      </div>
    {/if}

    <form onsubmit={handleSignUp} class="auth-form">
      <div class="form-group">
        <label for="username">Username</label>
        <input 
          type="text" 
          id="username" 
          bind:value={username} 
          required 
          autocomplete="username"
          placeholder="e.g. valid_user_name"
        />
      </div>

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

      <div class="form-group">
        <label for="password">Password</label>
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
        <label for="retype">Re-Type Password</label>
        <input 
          type="password" 
          id="retype" 
          bind:value={retype} 
          required 
          autocomplete="new-password"
          placeholder="••••••••"
        />
      </div>

      <button type="submit" class="btn btn-primary" disabled={loading || !!successMsg}>
        {loading ? 'Creating...' : 'Create Account'}
      </button>
    </form>

    <div class="auth-footer">
      <p>Already have an account? <a href="/user/login">Sign in</a></p>
    </div>
  </div>
</div>

<style>
  .auth-container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: calc(100vh - 100px);
    margin: 2rem 0;
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
    background-color: rgba(16, 185, 129, 0.1);
    border: 1px solid #10b981;
    color: #6ee7b7;
  }
</style>
