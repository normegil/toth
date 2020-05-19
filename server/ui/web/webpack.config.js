/* eslint-disable */
const path = require("path");
const VueLoaderPlugin = require("vue-loader/lib/plugin");
const HtmlWebpackPlugin = require("html-webpack-plugin");
const { CleanWebpackPlugin } = require("clean-webpack-plugin");

function getConfiguration(env) {
  let cfg = {
    entry: "./src/index.ts",
    mode: "development",
    optimization: {
      usedExports: true
    },
    devServer: {
      contentBase: './dist',
    },
    output: {
      filename: "bundle.js",
      path: path.resolve(__dirname, "dist")
    },
    module: {
      rules: [
        {
          test: /\.vue$/,
          loader: "vue-loader"
        },
        {
          test: /\.svg$/,
          loader: "vue-svg-loader" // `vue-svg` for webpack 1.x
        },
        {
          test: /\.(png|jpg|gif)$/,
          use: [
            {
              loader: "url-loader",
              options: {
                limit: 8192
              }
            }
          ]
        },
        {
          test: /\.tsx?$/,
          loader: "ts-loader",
          exclude: /node_modules/,
          options: {
            appendTsSuffixTo: [/\.vue$/]
          }
        },
        {
          test: /\.s?css$/,
          use: [
            "vue-style-loader",
            "css-loader",
            {
              loader: "sass-loader",
              options: {
                prependData: '@import "src/assets/scss/all";'
              }
            }
          ]
        },
        {
          test: /\.(woff(2)?|ttf|eot)(\?v=\d+\.\d+\.\d+)?$/,
          use: [
            {
              loader: "file-loader",
              options: {
                name: "[name].[ext]",
                outputPath: "fonts/"
              }
            }
          ]
        }
      ]
    },
    resolve: {
      extensions: [".tsx", ".ts", ".js", ".vue"],
      alias: {
        vue$: "vue/dist/vue.runtime.esm.js",
        "@": path.resolve("./src"),
        "@scss": path.resolve("./src/assets/scss")
      }
    },
    plugins: [
      new CleanWebpackPlugin(),
      new VueLoaderPlugin(),
      new HtmlWebpackPlugin({
        title: "Toth",
        template: "src/index.html"
      })
    ]
  };

  if (env === "dev") {
    cfg.mode = "development";
    cfg.optimization = {
      usedExports: true
    };
    return cfg;
  } else {
    cfg.mode = "production";
    return cfg;
  }
}

module.exports = getConfiguration;
