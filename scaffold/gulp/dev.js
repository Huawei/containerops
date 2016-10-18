'use strict';


let gulp = require('gulp');
let babel = require('gulp-babel');
let del = require('del');
let uglify = require('gulp-uglify');
// const rename = require('gulp-rename');
let concat = require('gulp-concat');
let browserify = require('browserify');
let source = require('vinyl-source-stream');
let sass = require('gulp-sass');
let imagemin = require('gulp-imagemin');
let gutil = require('gulp-util');
var replace = require('gulp-replace');
var browserSync = require('browser-sync').create();
var fs = require("fs");

gulp.task('dev:clean', function(){
     return del(['./dev']);
});

/**
 *  This will compile scss to css
 */
gulp.task('dev:styles', function () {
   return gulp.src('./src/sass/**/*.{scss,css}',{base:"./src/sass"})
        .pipe(sass().on('error', sass.logError))
        .pipe(gulp.dest('./dev/src/sass/'))
        .pipe(browserSync.stream())
        .on('error', gutil.log);
});


/**
 *  This will convert es6 to es5
 */
gulp.task('dev:babel', function(){
  return gulp.src('src/{app,vendor}/**/*.js',{base:"./src"})
    .pipe(babel({
      presets: ['es2015']
    }))
    // .pipe(uglify())
    .pipe(gulp.dest('dev/src'))
    .on('error', gutil.log);
})

/**
 *  This will copy templates to dev dist folder
 */
gulp.task('dev:html', function(){
  return gulp.src('src/**/*.html')
    // .pipe(uglify())
    .pipe(gulp.dest('dev/src'))
    .on('error', gutil.log);
})

/**
 *  This will copy assets to dev dist folder
 */
gulp.task('dev:images', function(){
  return gulp.src('src/assets/{images,svg}/*',{base:'./src/assets'})
    .pipe(imagemin())
    .pipe(gulp.dest('dev/src/assets'))
    .on('error', gutil.log);
})

gulp.task('dev:fonts', function(){
  return gulp.src('src/assets/fonts/**',{base:'./src/assets/fonts'})
    .pipe(gulp.dest('dev/src/assets/fonts'))
    .on('error', gutil.log);
})


/**
 *  This will browserify scripts
 */
gulp.task("dev:browserify",['dev:babel'],  () => {
    var b = browserify({
        entries: ["dev/src/app/index.js","dev/src/app/theme/settings.js","dev/src/app/theme/app.js"]
    });
    return b.bundle()
        .pipe(source("main.js"))
        .pipe(gulp.dest("dev/src"))
        .on('error', gutil.log);
});

/**
 *  This will concat all scripts include configed in scripts.json to one file: main.js
 */
gulp.task('dev:scripts',['dev:babel', 'dev:browserify'], function(done) {
  let config = JSON.parse(fs.readFileSync("src/scripts.json",'utf8'));
  let src = config.scripts.concat(['dev/src/main.js']);
  return gulp.src(src)
    .pipe(concat('main.js'))
    .pipe(gulp.dest('dev/src'))
    .on('error', gutil.log);
});

/**
 *  This will replace imported css in index.html
 */
gulp.task('dev:css-replace', ['dev:styles','dev:html'], function() {
  return gulp.src('dev/src/index.html')
      .pipe(replace(/<link rel="stylesheet">/g, '<link rel="stylesheet" href="sass/application.css" >'))
      .pipe(gulp.dest('dev/src'))
      .on('error', gutil.log);
});

/**
 *  This will replace imported script in index.html
 */
gulp.task('dev:script-replace', ['dev:scripts','dev:html'], function() {
  return gulp.src('dev/src/index.html')
      .pipe(replace(/<script\/>/g, '<script src="main.js"></script>'))
      .pipe(gulp.dest('dev/src'))
      .on('error', gutil.log);
});

gulp.task('dev:reload-js', ['dev:scripts'], function (done) {
    browserSync.reload();
    done();
});

gulp.task('dev:reload-html', ['dev:css-replace','dev:script-replace'], function () {
    browserSync.reload();
});
 /**
 *  This will watch files changing and do recompiling
 */
gulp.task("dev:watch",  () => {
   gulp.watch("./src/**/*.{scss,css}",['dev:styles']);
   gulp.watch(["src/{app,vendor}/**/*.js","src/scripts.json"], ['dev:reload-js']);
   gulp.watch("src/**/*.html",['dev:reload-html']);
   gulp.watch("src/assets/fonts/*",['dev:fonts','dev:styles']);
   gulp.watch("src/assets/images/*",['dev:images','dev:styles']);
   gulp.watch("src/assets/svg/*",['dev:images','dev:styles']);
});

/**
 *  This will start a server with browser-sync plugin
 */
 gulp.task('dev:browser-sync', ['dev:html','dev:images','dev:fonts','dev:css-replace','dev:script-replace'],() => {
   browserSync.init({
     server:{
        baseDir:"dev/src"
     }
   })
});


gulp.task('dev', ['dev:clean'], () => {
	 gulp.start('dev:html','dev:images','dev:fonts','dev:watch','dev:css-replace','dev:script-replace','dev:browser-sync');
});
