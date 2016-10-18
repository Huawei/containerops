'use strict';


// const gulp = require('gulp');
// const babel = require('gulp-babel');
// const uglify = require('gulp-uglify');
// const rename = require('gulp-rename');
// const concat = require('gulp-concat');
// const browserify = require('browserify');
// const source = require('vinyl-source-stream');
// const sass = require('gulp-sass');


// // styles
// gulp.task('styles', function () {

//     gulp.src('./sass/application.scss')
//         .pipe(sass().on('error', sass.logError))
//         .pipe(gulp.dest('./css'));

//     gulp.src('./sass/application.scss')
//         .pipe(sass({
//             outputStyle: 'compressed'
//         }).on('error', sass.logError))
//         .pipe(rename({suffix: '.min'}))
//         .pipe(gulp.dest('./css'));
// });




// // convert
// gulp.task('convertJS', function(){
//   return gulp.src('src/js/**/*.js',{base:"./src/js"})
//     .pipe(babel({
//       presets: ['es2015']
//     }))
//     // .pipe(uglify())
//     .pipe(gulp.dest('dist/js'))
// })



// // browserify
// gulp.task("browserify",['convertJS'], function () {
//     var b = browserify({
//         entries: "dist/js/index.js"
//     });

//     return b.bundle()
//         .pipe(source("main.js"))
//         .pipe(gulp.dest("dist/js"));
// });


// gulp.task('scripts',['browserify'], function() {
//   return gulp.src([
//     './bower_components/jquery/dist/jquery.min.js', 
//     './bower_components/jquery-pjax/jquery.pjax.js',
//     './bower_components/tether/dist/js/tether.js',
//     './bower_components/bootstrap/js/dist/util.js',
//     './bower_components/bootstrap/js/dist/collapse.js',
//     './bower_components/bootstrap/js/dist/tooltip.js',
//     './bower_components/bootstrap/js/dist/tab.js',
//     './bower_components/slimScroll/jquery.slimscroll.js',
//     './bower_components/widgster/widgster.js',
//     './node_modules/jsoneditor/dist/jsoneditor.min.js',
//     './dist/js/theme/settings.js',
//     './dist/js/theme/app.js',
//     './node_modules/d3/d3.min.js',  
//     './node_modules/node-uuid/uuid.js', 
//     './dist/js/main.js'
//     ])
//     .pipe(concat('app.js'))
//     .pipe(gulp.dest('./dist/js'));
// });

// gulp.task("watch",  () => {

//    gulp.watch("src/js/**/*.js",['scripts']);

// });

// gulp.task('default', ['convertJS','browserify','scripts','styles','watch']);



/** 
   new gulp content for changing directory
*/
let gulp = require('gulp');
let babel = require('gulp-babel');
let sass = require('gulp-sass');
let requireDir = require('require-dir');
let tasks = requireDir('./gulp');
 gulp.task('default', function () {
    gulp.start("dev")
 });


