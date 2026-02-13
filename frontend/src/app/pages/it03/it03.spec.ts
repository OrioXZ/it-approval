import { ComponentFixture, TestBed } from '@angular/core/testing';

import { It03 } from './it03';

describe('It03', () => {
  let component: It03;
  let fixture: ComponentFixture<It03>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [It03]
    })
    .compileComponents();

    fixture = TestBed.createComponent(It03);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
