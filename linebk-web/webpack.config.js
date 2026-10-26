module.exports = {
  module: {
    rules: [
      {
        test: /\.s[ac]ss$/i,
        type: 'css/auto',
        use: [
          {
            loader: 'sass-loader',
            options: {
              api: 'modern-compiler',
              sassOptions: {},
            },
          },
        ],
      },
    ],
  },
  experiments: {
    css: true,
  },
};
