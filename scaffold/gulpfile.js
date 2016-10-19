'use strict';


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


