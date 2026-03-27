function jsonHeaders(extraHeaders = {}) {
  return { "Content-Type": "application/json", ...extraHeaders };
}

export function authHeaders(token, extraHeaders = {}) {
  if (!token) {
    return jsonHeaders(extraHeaders);
  }

  return jsonHeaders({
    Authorization: `Bearer ${token}`,
    ...extraHeaders,
  });
}

export async function requestJSON(url, options = {}) {
  const headers = options.isFormData ? (options.headers || {}) : jsonHeaders(options.headers || {});
  const response = await fetch(url, {
    method: options.method || "GET",
    headers,
    body: options.body,
  });

  const contentType = response.headers.get("content-type") || "";
  const payload = contentType.includes("application/json")
    ? await response.json()
    : { message: await response.text() };

  if (!response.ok) {
    throw new Error(payload.message || payload.error || `请求失败：${response.status}`);
  }

  return payload;
}

export function fetchTracks() {
  return requestJSON("/music/local/list");
}

export function uploadTracks(files) {
  const formData = new FormData();
  files.forEach((file) => formData.append("files", file));
  return requestJSON("/music/local/upload", {
    method: "POST",
    body: formData,
    isFormData: true,
  });
}

export function loginUser(credentials) {
  return requestJSON("/user/login", {
    method: "POST",
    body: JSON.stringify(credentials),
  });
}

export function registerUser(payload) {
  return requestJSON("/user/register", {
    method: "POST",
    body: JSON.stringify({ user_info: payload }),
  });
}

export function fetchCurrentUser(token) {
  return requestJSON("/user/me", {
    headers: authHeaders(token),
  });
}

export function toggleLike(token, fileName) {
  return requestJSON("/user/likes/toggle", {
    method: "POST",
    headers: authHeaders(token),
    body: JSON.stringify({ file_name: fileName }),
  });
}

export function removeToken(username, accessToken) {
  return requestJSON("/token/remove", {
    method: "POST",
    body: JSON.stringify({
      username,
      access_token: accessToken,
    }),
  });
}
