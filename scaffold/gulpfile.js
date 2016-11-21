'use strict';


/** 
   new gulp content for changing directory
*/
var gulp = require('gulp');
var babel = require('gulp-babel');
var sass = require('gulp-sass');
var requireDir = require('require-dir');
var tasks = requireDir('./gulp');
 gulp.task('default', function () {
    gulp.start("dev")
 });


