module.exports = {
  presets: [
    ['@babel/preset-env', {targets: {esmodules: true, node: 'current'}}],
    [ '@babel/preset-react', { runtime: 'automatic' } ],
    '@babel/preset-typescript',
  ],
  plugins: [
    [
      'module-resolver', {
        root: ["/app"],
        alias: {
          "@/fonts": "./app/styles/fonts",
          "@/app/ui/dashboard": "./app/ui/dashboard",
        },
      },
    ]
  ]
};
