'use strict';


/** 
   new gulp content for changing directory
*/
var gulp = require('gulp');
var sass = require('gulp-sass');
var requireDir = require('require-dir');
var tasks = requireDir('./gulp');
 gulp.task('default', function () {
    gulp.start("dev")
 });


