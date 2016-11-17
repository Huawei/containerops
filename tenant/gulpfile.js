'use strict';


let gulp = require('gulp');
let sass = require('gulp-sass');


/**
 *  This will compile scss to css
 */
gulp.task('default', function() {
    return gulp.src('./sass/**/*.{scss,css}', { base: "./sass" })
        .pipe(sass().on('error', sass.logError))
        .pipe(gulp.dest('./css/'));
});


