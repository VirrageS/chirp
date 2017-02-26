let env = {
  production: false,
  api_url: "http://0.0.0.0:8080"
}
if (process.env.ENV === "production") {
  env = {
    production: true,
    api_url: "http://backend.show"
  };
}

export const environment = env;