/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

'use strict';


var gulp = require('gulp');
var babel = require('gulp-babel');
var del = require('del');
var uglify = require('gulp-uglify');
// const rename = require('gulp-rename');
var concat = require('gulp-concat');
var browserify = require('browserify');
var source = require('vinyl-source-stream');
var sass = require('gulp-sass');
var imagemin = require('gulp-imagemin');
var gutil = require('gulp-util');
var replace = require('gulp-replace');
var browserSync = require('browser-sync').create();
var minimist = require('minimist');
var fs = require("fs");


var args = minimist(process.argv.slice(2), {
    boolean: ["dist"]
});

gulp.task('dev:clean', function() {
    return del(['./dev']);
});

/**
 *  This will compile scss to css
 */
gulp.task('dev:styles', function() {
    return gulp.src('./src/sass/**/*.{scss,css}', { base: "./src/sass" })
        .pipe(sass().on('error', sass.logError))
        .pipe(gulp.dest('./dev/src/sass/'))
        .pipe(browserSync.stream())
        .on('error', gutil.log);
});


/**
 *  This will convert es6 to es5
 */
gulp.task('dev:babel', function() {
    return gulp.src('src/{app,vendor}/**/*.js', { base: "./src" })
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
gulp.task('dev:html', function() {
    return gulp.src('src/**/*.html')
        // .pipe(uglify())
        .pipe(gulp.dest('dev/src'))
        .on('error', gutil.log);
})

/**
 *  This will copy assets to dev dist folder
 */
gulp.task('dev:images', function() {
    return gulp.src('src/assets/{images,svg}/*', { base: './src/assets' })
        .pipe(imagemin())
        .pipe(gulp.dest('dev/src/assets'))
        .on('error', gutil.log);
})

gulp.task('dev:fonts', function() {
    return gulp.src('src/assets/fonts/**', { base: './src/assets/fonts' })
        .pipe(gulp.dest('dev/src/assets/fonts'))
        .on('error', gutil.log);
})

/**
 *  This will copy host json to dev dist folder
 */
gulp.task('dev:json', function() {
    return gulp.src('src/host.json')
        // .pipe(uglify())
        .pipe(gulp.dest('dev/src'))
        .on('error', gutil.log);
})

/**
 *  This will browserify scripts
 */
gulp.task("dev:browserify", ['dev:babel'], function() {
    var b = browserify({
        entries: ["dev/src/app/index.js", "dev/src/app/theme/settings.js", "dev/src/app/theme/app.js","dev/src/app/history/paginate.js" ]
    });
    return b.bundle()
        .pipe(source("main.js"))
        .pipe(gulp.dest("dev/src"))
        .on('error', gutil.log);
});

/**
 *  This will concat all scripts include configed in scripts.json to one file: main.js
 */
gulp.task('dev:scripts', ['dev:babel', 'dev:browserify'], function(done) {
    let config = JSON.parse(fs.readFileSync("src/scripts.json", 'utf8'));
    let src = config.scripts.concat(['dev/src/main.js']);
    return gulp.src(src)
        .pipe(concat('main.js'))
        .pipe(gulp.dest('dev/src'))
        .on('error', gutil.log);
});

/**
 *  This will replace imported css in index.html
 */
gulp.task('dev:css-replace', ['dev:styles', 'dev:html'], function() {
    return gulp.src('dev/src/index.html')
        .pipe(replace(/<link rel="stylesheet">/g, '<link rel="stylesheet" href="sass/application.css" >'))
        .pipe(gulp.dest('dev/src'))
        .on('error', gutil.log);
});

/**
 *  This will replace imported script in index.html
 */
gulp.task('dev:script-replace', ['dev:scripts', 'dev:html'], function() {
    var randomCopy = "?copy=" + Math.random();
    return gulp.src('dev/src/index.html')
        .pipe(replace(/<script\/>/g, '<script src="main.js'+randomCopy+'"></script>'))
        .pipe(gulp.dest('dev/src'))
        .on('error', gutil.log);
});

gulp.task('dev:reload-js', ['dev:scripts'], function(done) {
    browserSync.reload();
    done();
});

gulp.task('dev:reload-html', ['dev:css-replace', 'dev:script-replace'], function() {
    browserSync.reload();
});
/**
 *  This will watch files changing and do recompiling
 */
gulp.task("dev:watch", function() {
    gulp.watch("./src/**/*.{scss,css}", ['dev:styles']);
    gulp.watch(["src/{app,vendor}/**/*.js", "src/scripts.json"], ['dev:reload-js']);
    gulp.watch("src/**/*.html", ['dev:reload-html']);
    gulp.watch("src/assets/fonts/*", ['dev:fonts', 'dev:styles']);
    gulp.watch("src/assets/images/*", ['dev:images', 'dev:styles']);
    gulp.watch("src/assets/svg/*", ['dev:images', 'dev:styles']);
});

/**
 *  This will start a server with browser-sync plugin
 */
gulp.task('dev:browser-sync', ['dev:html', 'dev:images', 'dev:fonts', 'dev:css-replace', 'dev:script-replace'], function() {
    browserSync.init({
        server: {
            baseDir: "dev/src"
        },
        ghostMode: false
    })
});

gulp.task('dev:copy', ['dev:html', 'dev:images', 'dev:fonts', 'dev:css-replace', 'dev:script-replace'], function() {
    if (args.dist) {
        return gulp.src('dev/**')
            .pipe(gulp.dest(args._[0]))
            .on('error', gutil.log);
    }
});

gulp.task('dev', ['dev:clean'], function() {
    if (args.dist) {
        gulp.start('dev:html', 'dev:images', 'dev:fonts','dev:json', 'dev:css-replace', 'dev:script-replace', 'dev:copy');
    } else {
        gulp.start('dev:html', 'dev:images', 'dev:fonts','dev:json', 'dev:watch', 'dev:css-replace', 'dev:script-replace', 'dev:browser-sync');
    }

});
