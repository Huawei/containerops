'use strict';


let gulp = require('gulp');
let sass = require('gulp-sass');


/**
 *  This will compile scss to css
 */
gulp.task("dev:watch", () => {
    gulp.watch("./sass/**/*.{scss,css}", ['dev:style']);
});
gulp.task("dev:style", () => {
	return gulp.src('./sass/**/*.{scss,css}', { base: "./sass" })
        .pipe(sass().on('error', sass.logError))
        .pipe(gulp.dest('./css/'));
})
gulp.task("default", () => {
	gulp.start('dev:watch','dev:style');
})


